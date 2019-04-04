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

	"github.com/gobwas/glob"
	"github.com/mitchellh/go-homedir"
	"github.com/moisespsena-go/error-wrap"
)

type ValidFunc func(pth string, info os.FileInfo) bool

type Patterns struct {
	values []glob.Glob
	m      map[string]bool
}

func (e *Patterns) Append(values ...string) (err error) {
	var g glob.Glob
	if e.m == nil {
		e.m = map[string]bool{}
	}

	for i, value := range values {
		if len(value) == 0 {
			continue
		}
		if _, ok := e.m[value]; ok {
			continue
		}
		if g, err = PathGlobCompile(value); err != nil {
			return errwrap.Wrap(err, "Value %d: %q", i, value)
		}

		e.m[value] = true
		e.values = append(e.values, g)
	}
	return nil
}

func (e *Patterns) Values() []glob.Glob {
	return e.values
}

func (e *Patterns) ValidFunc() ValidFunc {
	return func(pth string, info os.FileInfo) bool {
		if len(e.values) == 0 {
			return true
		}
		for _, g := range e.values {
			if g.Match(pth) {
				return true
			}
		}
		return false
	}
}

func (e *Patterns) ExcludeFunc() ValidFunc {
	return func(pth string, info os.FileInfo) bool {
		if len(e.values) == 0 {
			return false
		}
		for _, g := range e.values {
			if g.Match(pth) {
				return true
			}
		}
		return false
	}
}

func Exclude(values ...string) (ValidFunc, error) {
	p := &Patterns{}
	if err := p.Append(values...); err != nil {
		return nil, err
	}
	if len(p.values) == 0 {
		return nil, nil
	}
	return p.ExcludeFunc(), nil
}

type BackupOptions struct {
	Target        string
	DefaultBackup bool
	Writer        io.Writer
	Patterns      Patterns
}

func (env *GoEnvCmd) Backup(name string, options *BackupOptions) error {
	pth, err := env.Env.Backup(name, options)
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

	excludeFile := filepath.Join(pth, ".goenv_settings", "backup_exclude")
	ok, err := IsFile(excludeFile)
	if err != nil {
		return "", errwrap.Wrap(err, "Check file %q", excludeFile)
	}
	if ok {
		exclude, err := readLines(excludeFile)
		if err != nil {
			return "", err
		}
		err = options.Patterns.Append(exclude...)
		if err != nil {
			return "", errwrap.Wrap(err, "Parse default exclude patterns in %q", excludeFile)
		}
	}

	if options.Target != "" {
		writer, err := os.Create(options.Target)
		if err != nil {
			return "", err
		}
		defer writer.Close()
		err = compress(pth, writer, options.Patterns.ExcludeFunc())
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
		err = compress(pth, writer, options.Patterns.ExcludeFunc())
		return target, err
	}

	if options.Writer != nil {
		return "", compress(pth, options.Writer, options.Patterns.ExcludeFunc())
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
	pth, err := env.Env.Restore(options)
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
