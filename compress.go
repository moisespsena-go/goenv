// Copyright © 2018 Moises P. Sena <moisespsena@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goenv

import (
	"time"
	"archive/tar"
	"os"
	"path/filepath"
	"strings"
	"io"
	"compress/gzip"
	"fmt"
	"github.com/dustin/go-humanize"
)

const (
	Verbose ExtractOptions = 1 << iota
	Trial
)

type ExtractOptions int

func (o ExtractOptions) IsVerbose() bool {
	return (o & Verbose) != 0
}

func (o ExtractOptions) IsTrial() bool {
	return (o & Trial) != 0
}

// The gzip file stores a header giving metadata about the compressed file.
// That header is exposed as the fields of the Writer and Reader structs.
type Header struct {
	Comment string    // comment
	Extra   []byte    // "extra data"
	ModTime time.Time // modification time
	Name    string    // file name
	OS      byte      // operating system type
}

func compress(source string, writer io.Writer) error {
	gzWriter := gzip.NewWriter(writer)
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()
	defer gzWriter.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarWriter, file)
			return err
		})
}

type BackupFile struct {
	Reader *tar.Reader
	first  *tar.Header
}

func NewBackupReader(reader io.Reader, archive bool) (bkp *BackupFile, err error) {
	if !archive {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
	}
	bkp = &BackupFile{Reader: tar.NewReader(reader)}
	return
}

func (b *BackupFile) GetRootName() (name string, err error) {
	if b.first != nil {
		return b.first.Name, nil
	}
	err = b.Each(func(header *tar.Header, reader *tar.Reader) error {
		b.first = header
		return io.EOF
	})

	if err != nil {
		return "", err
	}

	name = strings.Trim(b.first.Name, string(os.PathSeparator))
	if !b.first.FileInfo().IsDir() || strings.Contains(name, string(os.PathSeparator)) {
		return "", fmt.Errorf("Invalid root name %q", name)
	}

	return name, nil
}

func (b *BackupFile) Each(cb func(header *tar.Header, reader *tar.Reader) error) error {
	for {
		header, err := b.Reader.Next()

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		err = cb(header, b.Reader)

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

func (b *BackupFile) EachRoot(rootName string, cb func(header *tar.Header, reader *tar.Reader) error) error {
	originalRootName, err := b.GetRootName()
	if err != nil {
		return err
	}
	if rootName == "" || rootName == originalRootName {
		err = cb(b.first, b.Reader)
		if err != nil {
			return err
		}
		return b.Each(cb)
	}

	lorn := len(originalRootName)
	header := *b.first
	header.Name = rootName + header.Name[lorn:]
	err = cb(&header, b.Reader)
	if err != nil {
		return err
	}

	return b.Each(func(header *tar.Header, reader *tar.Reader) error {
		header.Name = filepath.Join(rootName, header.Name[lorn:])
		return cb(header, reader)
	})
}

func pad(v string, count int) string {
	r := strings.Repeat(" ", count-len(v))
	return v + r
}
func padl(v string, count int) string {
	r := strings.Repeat(" ", count-len(v))
	return r + v
}

var typeDesc = map[byte]string{
	tar.TypeReg:           "F",
	tar.TypeRegA:          "Fa",
	tar.TypeLink:          "L",
	tar.TypeSymlink:       "S",
	tar.TypeChar:          "C",
	tar.TypeBlock:         "B",
	tar.TypeDir:           "D",
	tar.TypeFifo:          "FF",
	tar.TypeCont:          "Cn",
	tar.TypeXHeader:       "Xh",
	tar.TypeXGlobalHeader: "XH",
	tar.TypeGNUSparse:     "gnuS",
	tar.TypeGNULongName:   "gnuL",
	tar.TypeGNULongLink:   "gnuK",
}

func (b *BackupFile) Extract(rootName, target string, options ExtractOptions) error {
	return b.EachRoot(rootName, func(header *tar.Header, reader *tar.Reader) (err error) {
		info := header.FileInfo()
		if options.IsVerbose() {
			prefix := pad(typeDesc[header.Typeflag], 5) + " "
			switch header.Typeflag {
			case tar.TypeReg,tar.TypeRegA:
				prefix += pad("["+humanize.Bytes(uint64(info.Size()))+"]", 12)
			default:
				prefix += pad("", 12)
			}
			name := header.Name
			switch header.Typeflag {
			case tar.TypeLink, tar.TypeSymlink:
				name += " -> " + header.Linkname
			}
			os.Stdout.WriteString(prefix + name + "... ")
		}
		path := filepath.Join(target, header.Name)
		if info.IsDir() {
			if options.IsVerbose() {
				defer func() {
					if err != nil {
						os.Stdout.WriteString("failed.\n")
					} else {
						os.Stdout.WriteString("done.\n")
					}
				}()
			}
			if options.IsTrial() {
				return nil
			}
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			return nil
		}

		if options.IsTrial() {
			os.Stdout.WriteString("done.\n")
			return
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}

		defer func() {
			if err == nil {
				err = file.Close()
			} else {
				file.Close()
			}
			if options.IsVerbose() {
				if err != nil {
					os.Stdout.WriteString("failed.\n")
				} else {
					os.Stdout.WriteString("done.\n")
				}
			}
		}()
		_, err = io.Copy(file, reader)

		return err
	})
}
