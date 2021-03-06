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
	"github.com/moisespsena-go/error-wrap"
	"github.com/moisespsena/go-ioutil"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func IsDir(pth ...string) (bool, error) {
	p := filepath.Join(pth...)
	s, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("'%v': %v", p, err)
	}
	if !s.IsDir() {
		return false, fmt.Errorf("'%v': Isn't directory.", p)
	}
	return true, nil
}

func IsFile(pth ...string) (bool, error) {
	p := filepath.Join(pth...)
	s, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("'%v': %v", p, err)
	}
	if s.IsDir() {
		return false, fmt.Errorf("'%v': Is directory.", p)
	}
	return true, nil
}

func MkdirAll(pth ...string) error {
	p := filepath.Join(pth...)
	parent := p
	mode := os.FileMode(0777)
	var (
		err  error
		stat os.FileInfo
	)
	old := ""
	for old != parent {
		stat, err = os.Stat(parent)
		if err != nil {
			if os.IsNotExist(err) {
				old = parent
				parent = filepath.Dir(parent)
				continue
			}
			return fmt.Errorf("'%v': %v", parent, err)
		}
		mode = stat.Mode()
		err = os.MkdirAll(p, mode)
		if err != nil {
			return fmt.Errorf("'%v': %v", p, err)
		}
		return nil
		break
	}
	return fmt.Errorf("'%v': Invalid path.", p)
}

func TimeString(t time.Time) string {
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d%v", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(),
		t.Second(), t.Nanosecond())
}

func readLines(pth string) ([]string, error) {
	f, err := os.Open(pth)
	if err != nil {
		return nil, errwrap.Wrap(err, "Open %q", pth)
	}
	defer f.Close()

	var (
		lines []string
		line  []byte
	)
	i := 0
	for err == nil {
		i++
		line, err = iolr.ReadLine(f)
		if err == nil {
			if line[len(line)-1] == '\r' {
				line = line[0 : len(line)-1]
			}
			lines = append(lines, strings.TrimSpace(string(line)))
		}
	}
	if err != nil && err != io.EOF {
		return nil, errwrap.Wrap(err, "Read Line %d", i)
	}
	return lines, nil
}

const activateData = `
export GOPATH="$GOENVROOT/$GOENVNAME"
export OLDPS1=$PS1
export PS1="[go:$GOENVNAME] $PS1"
export OLDPATH="$PATH"
export PATH="$GOPATH/bin:$PATH"
alias gcd="cd $GOPATH"
goenv-deactivate() {
	export PS1=$OLDPS1
	export PATH=$OLDPATH
	unset GOPATH
	unset OLDPS1
	unset OLDPATH
	unalias gcd
	unset -f goenv-deactivate
}
`
