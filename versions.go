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
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"

	"os"
	"path/filepath"

	"github.com/cavaliercoder/grab"
	"github.com/dustin/go-humanize"
	errwrap "github.com/moisespsena-go/error-wrap"
)

const VERSIONS_BASENAME = ".goversions"

var versionRe, _ = regexp.Compile(`\D(\d+)`)
var versionRe2, _ = regexp.Compile(`^(\D+)\d.*$`)

type GoBinVersion struct {
	Version string
	OsInfo  string
}

func NewGoBinVersion(binName string) (*GoBinVersion, error) {
	ew := func(child error, self interface{}, args ...interface{}) error {
		return errwrap.Wrap(errwrap.Wrap(child, self, args...), "GoBinVersion of %q", binName)
	}
	cmd := exec.Command(binName, "version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Start()
	if err != nil {
		return nil, ew(err, "Cmd Start")
	}
	err = cmd.Wait()
	if err != nil {
		return nil, ew(err, "Cmd Wait")
	}
	parts := strings.Split(strings.TrimSpace(out.String()), " ")
	osinfo := strings.Replace(parts[len(parts)-1], "/", "-", 1)
	return &GoBinVersion{parts[2], osinfo}, nil
}

func GetSystemGoVersion() (*GoVersion, error) {
	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		var pth string
		for _, p := range strings.Split(os.Getenv("PATH"), string(filepath.ListSeparator)) {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			pth = filepath.Join(p, "go")
			_, err := os.Stat(pth)
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return nil, errwrap.Wrap(err, "Get Stat of %q", pth)
			}
			absPath, err := filepath.Abs(p)
			if err != nil {
				return nil, errwrap.Wrap(err, "Get Abs Path of %q", p)
			}
			goroot = filepath.Dir(absPath)
			break
		}
	}

	if goroot != "" {
		v, err := NewGoVersion(goroot)
		if v != nil {
			v.System = true
		}
		return v, err
	}
	return nil, nil
}

type GoVersion struct {
	versions     *GoVersions
	ID           string
	Name         string
	UpdatedAt    time.Time
	downloadUrl  string
	Installed    bool
	Root         string
	System       bool
	BinVersion   *GoBinVersion
	downloadPath string
}

func NewGoVersion(goroot string) (v *GoVersion, err error) {
	pth := filepath.Join(goroot, "bin", "go")
	binVersion, err := NewGoBinVersion(pth)
	if err != nil {
		return nil, err
	}
	return &GoVersion{Root: goroot, Name: binVersion.Version, BinVersion: binVersion}, nil
}
func (v *GoVersion) DownloadPath() string {
	if v.downloadPath == "" {
		v.downloadPath = filepath.Join(v.versions.Dir(), filepath.Base(v.DownloadUrl()))
	}
	return v.downloadPath
}
func (v *GoVersion) Downloadable(client *http.Client) (bool, error) {
	r, err := client.Head(v.DownloadUrl())
	defer r.Body.Close()
	if err == nil && r.StatusCode == 200 {
		return true, nil
	}
	return false, errwrap.Wrap(err, "HTTP HEAD %q", v.downloadUrl)
}

func (v *GoVersion) DownloadUrl() string {
	return v.downloadUrl
}

type GoVersions struct {
	Env *GoEnv
}

func NewGoVersions(env *GoEnv) *GoVersions {
	return &GoVersions{env}
}

func (v *GoVersions) Dir() string {
	return filepath.Join(v.Env.DbDir, VERSIONS_BASENAME)
}

func (v *GoVersions) DirExists() (pth string, exists bool, err error) {
	pth = v.Dir()
	s, err := os.Stat(pth)
	if err != nil {
		if os.IsNotExist(err) {
			return pth, false, nil
		}
		return "", false, errwrap.Wrap(err, "file stat")
	}
	if !s.IsDir() {
		return "", false, fmt.Errorf("%q is not a directory.", pth)
	}
	return pth, true, err
}

func (v *GoVersions) Ls() (vs []*GoVersion, err error) {
	dir, exists, err := v.DirExists()
	if err != nil {
		return nil, err
	}
	if exists {
		s, err := os.Open(dir)
		if err != nil {
			return nil, errwrap.Wrap(err, "Open %q", dir)
		}
		items, err := s.Readdir(-1)
		if err != nil {
			return nil, errwrap.Wrap(err, "Readdir %q", dir)
		}
		for _, f := range items {
			if f.IsDir() {
				pth := filepath.Join(dir, f.Name())
				version, err := NewGoVersion(pth)
				if version != nil {
					version.Name = f.Name()
					version.versions = v
					if err != nil {
						return nil, errwrap.Wrap(err, "Parse Version of %q", pth)
					}
					vs = append(vs, version)
				}
			}
		}
	}
	return
}

func (v *GoVersions) Set(versionName, envName string) (err error) {
	system, err := GetSystemGoVersion()
	if err != nil {
		return err
	}
	if system == nil {
		return errwrap.Wrap(err, "GO isn't available on system. Please install it from https://golang.org/dl")
	}

	if versionName == "sys" {
		err = v.Env.SetGoVersion(versionName, "")
	} else {
		versionName = strings.ToLower(versionName)
		versions, err := v.Ls()
		if err != nil {
			return errwrap.Wrap(err, "Get installed versions")
		}

		var version *GoVersion

		for _, ver := range versions {
			if ver.Name == versionName {
				version = ver
				break
			}
		}

		if version == nil {
			return fmt.Errorf("GoLang version %q has not be installed", versionName)
		}

		err = v.Env.SetGoVersion(envName, filepath.Join("$GOENVROOT", VERSIONS_BASENAME, version.Name))
	}
	return errwrap.Wrap(err, "Env Set Go Version")
}

func (v *GoVersions) Download(names ...string) (versions []*GoVersion, err error) {
	if len(names) == 0 {
		return
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return nil, errwrap.Wrap(err, "GOPATH enviromente variable isn't defined.")
	}

	dir, exists, err := v.DirExists()
	if err != nil {
		return nil, err
	}

	versions, err = v.Available(names...)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return nil, errwrap.Wrap(err, "Create versions directory.")
		}
	}

	// create client
	client := grab.NewClient()

	l := len(versions)
	resp := make([]*grab.Response, len(versions))

	hb := func(value int64) string {
		return humanize.Bytes(uint64(value)) + " (" + strconv.Itoa(int(value)) + " bytes)"
	}

	for i, v := range versions {
		v.Root = filepath.Join(dir, v.Name)
		req, _ := grab.NewRequest(v.DownloadPath(), v.DownloadUrl())
		fmt.Printf("[%v] Downloading %v... ", v.Name, req.URL())
		resp[i] = client.Do(req)
		fmt.Printf("[HTTP %v]", resp[i].HTTPResponse.Status)
		switch resp[i].HTTPResponse.StatusCode {
		case 200, 206:
			fmt.Printf(" Left Size: %v", hb(resp[i].HTTPResponse.ContentLength))
		}
		fmt.Println()
	}

	var dok []*GoVersion
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()

	for l > 0 {
		<-t.C
		for i, r := range resp {
			if r == nil {
				continue
			}
			select {
			case <-r.Done:
				resp[i] = nil
				l--

				// check for errors
				if err := r.Err(); err != nil {
					fmt.Fprintf(os.Stderr, "[%v] Download failed: %v\n", versions[i].Name, err)
				} else {
					dok = append(dok, versions[i])
					fmt.Printf("[%v] Download saved to %v\n", versions[i].Name, r.Filename)
				}
			default:
				fmt.Printf("[%v] transferred %v / %v (%.2f%%)\n",
					versions[i].Name, hb(r.BytesComplete()), hb(r.Size), 100*r.Progress())
			}
		}
	}
	return dok, nil
}

func (vs *GoVersions) Install(names ...string) (versions []*GoVersion, err error) {
	versions, err = vs.Download(names...)
	if err != nil {
		return
	}
	for _, v := range versions {
		_, err = os.Stat(v.Root)
		if err == nil || !os.IsNotExist(err) {
			err = os.RemoveAll(v.Root)
			if err != nil {
				return nil, errwrap.Wrap(err, "Remove %q", v.Root)
			}
		} else if !os.IsNotExist(err) {
			return nil, errwrap.Wrap(err, "Stat of %q", v.Root)
		}
		f, err := os.Open(v.DownloadPath())
		if err != nil {
			return nil, errwrap.Wrap(err, "Open %q", v.downloadPath)
		}
		bkp, err := NewBackupReader(f, false)
		if err != nil {
			return nil, errwrap.Wrap(err, "Reader %q", v.downloadPath)
		}
		fmt.Printf("[%v] Extract %q to %q...", v.Name, v.downloadPath, v.Root)

		err = bkp.Extract(v.Name, vs.Dir(), ExtractOptions(0))
		fmt.Println(" Done.")
		if err != nil {
			return nil, errwrap.Wrap(err, "Extract %q", v.downloadPath)
		}
	}
	return
}

func (v *GoVersions) Available(terms ...string) ([]*GoVersion, error) {
	system, err := GetSystemGoVersion()
	if err != nil {
		return nil, err
	}
	if system == nil {
		return nil, errwrap.Wrap(err, "GO isn't available on system. Please install it from https://golang.org/dl")
	}

	r, err := http.Get("https://golang.org/dl")
	if err != nil {
		return nil, errwrap.Wrap(err, "Get Milestones from %q", r.Request.URL)
	}
	defer r.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return nil, errwrap.Wrap(err, "Decode go download page failed")
	}

	var versions []*GoVersion

	var idre, _ = regexp.Compile(`^go\d+\.`)
	// Find the review items
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("id"); idre.MatchString(name) {
			s.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
				fileName := s.Find("td").Eq(0)
				fileNameS := fileName.Text()
				if !strings.HasSuffix(fileNameS, ".tar.gz") {
					return
				}
				fileNameS = strings.TrimSuffix(fileNameS, ".tar.gz")
				parts := strings.Split(strings.TrimPrefix(fileNameS, name+"."), "-")
				if len(parts) != 2 {
					return
				}
				if parts[0] != runtime.GOOS {
					return
				}
				if parts[1] != runtime.GOARCH {
					return
				}
				url, _ := fileName.Find("a").Attr("href")
				if url == "" {
					return
				}

				ver := &GoVersion{
					Name:        name,
					ID:          sname(name),
					downloadUrl: url,
					versions:    v,
				}

				root := filepath.Join(v.Dir(), ver.Name)
				if _, err := os.Stat(root); err == nil {
					ver.Root = root
				}

				versions = append(versions, ver)
			})
		}
	})

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].ID > versions[j].ID
	})

	return versions, nil
}

func sname(name string) (v string) {
	parts := strings.Split(name[2:], ".")

	for i, s := range parts {
		for j, x := range s {
			if !unicode.IsDigit(x) {
				s = s[0:j]
				break
			}
		}
		si, _ := strconv.Atoi(s)
		parts[i] = fmt.Sprintf("%04d", si)
	}
	return strings.Join(parts, "")
}
