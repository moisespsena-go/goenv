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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
)

type BackupOptions struct {
	Target        string
	DefaultBackup bool
	Writer        io.Writer
}

func (env *GoEnvCmd) Backup(name string, options *BackupOptions) error {
	pth, err := env.env.Backup(name, options)
	if err != nil {
		return fmt.Errorf("Backup for %q failed: %v", name, err)
	}
	if pth != "" {
		fmt.Fprintf(os.Stdout, "Backup save on %q\n", pth)
	}
	return nil
}

func (env *GoEnv) TempDir() (pth string, err error) {
	pth = filepath.Join(env.DbDir, ".tmp")
	exists, err := IsDir(pth)

	if err != nil {
		return "", err
	}

	if !exists {
		if err = MkdirAll(pth); err != nil {
			return "", err
		}
	}
	return
}

func (env *GoEnv) Backup(name string, options *BackupOptions) (string, error) {
	pth, err := env.GetCheck(name)

	if err != nil {
		return "", err
	}

	if options.Target != "" {
		writer, err := os.Create(options.Target)
		if err != nil {
			return "", err
		}
		defer writer.Close()
		err = compress(pth, writer)
		return options.Target, err
	}

	if options.DefaultBackup {
		bkpDir := filepath.Join(env.DbDir, ".backup", name)
		exists, err := IsDir(bkpDir)

		if err != nil {
			return "", err
		}

		if !exists {
			if err = MkdirAll(bkpDir); err != nil {
				return "", err
			}
		}

		target := filepath.Join(bkpDir, name+"_"+TimeString(time.Now())+".tar.gz")
		writer, err := os.Create(target)
		if err != nil {
			return "", err
		}
		defer writer.Close()
		err = compress(pth, writer)
		return target, err
	}

	if options.Writer != nil {
		return "", compress(pth, options.Writer)
	}

	return "", fmt.Errorf("No target defined.")
}

type RestoreOptions struct {
	Source    string
	OverWrite bool
	Update    bool
	Reader    io.ReadSeeker
	Name      string
	Archive   bool
	Verbose   bool
	Trial     bool
}

func (env *GoEnvCmd) Restore(options *RestoreOptions) error {
	pth, err := env.env.Restore(options)
	if err != nil {
		return fmt.Errorf("Restore failed: %v", err)
	}
	if pth != "" {
		fmt.Fprintf(os.Stdout, "Restore saved on %q\n", pth)
	}
	return nil
}

func (env *GoEnv) Restore(options *RestoreOptions) (string, error) {
	var err error
	var reader io.ReadSeeker

	if options.Source != "" {
		src, err := homedir.Expand(options.Source)
		if err != nil {
			return "", err
		}
		reader, err = os.Open(src)
		if err != nil {
			return "", err
		}
		defer reader.(io.ReadCloser).Close()
	} else if options.Reader != nil {
		reader = options.Reader
	}

	if reader == nil {
		return "", fmt.Errorf("No source defined.")
	}

	//defer os.RemoveAll(dir)
	bkp, err := NewBackupReader(reader, options.Archive)
	if err != nil {
		return "", err
	}
	name, err := bkp.GetRootName()
	if err != nil {
		return "", err
	}

	if options.Name != "" {
		name = options.Name
	}

	pth, err := env.GetPath(name, false)

	if err != nil {
		return "", err
	}

	exists, err := IsDir(pth)
	if err != nil {
		return "", err
	}

	if !exists || (options.Update || options.OverWrite) {
		if exists && options.OverWrite && !options.Trial {
			err := os.RemoveAll(pth)
			if err != nil {
				return pth, err
			}
		}
		opts := ExtractOptions(0)
		if options.Verbose {
			opts |= Verbose
		}
		if options.Trial {
			opts |= Trial
		}
		err = bkp.Extract(name, env.DbDir, opts)
		if err != nil {
			return "", err
		}
		return pth, nil
	}

	return pth, fmt.Errorf("Enviroment %q on %q exists.", name, pth)
}
