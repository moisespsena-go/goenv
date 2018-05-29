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
	"path/filepath"
	"io/ioutil"
	"os"
	"time"
)

type GoEnv struct {
	DbDir string
}

func NewGoEnv(dbDir string, check bool) (env *GoEnv, err error) {
	ok, err := IsDir(dbDir)
	if err != nil {
		return nil, err
	}
	if check && !ok {
		return nil, fmt.Errorf("Database %q isn't initialized.", dbDir)
	}
	return &GoEnv{dbDir}, nil
}

func (env *GoEnv) Init(name string) (err error) {
	var ok bool
	pth := filepath.Join(env.DbDir, name)
	ok, err = IsDir(pth, "src")
	if err != nil {
		return err
	}
	if !ok {
		err = MkdirAll(pth, "src")
		if err != nil {
			return err
		}
	}
	err = MkdirAll(pth, "bin")
	if err != nil {
		return err
	}

	err = createActivate(pth)
	if err != nil {
		return err
	}
	return nil
}

func (env *GoEnv) Ls() (names []string, err error) {
	files, err := ioutil.ReadDir(env.DbDir)
	if err != nil {
		return nil, fmt.Errorf("'%v': %v", env.DbDir, err)
	}
	for _, f := range files {
		if f.IsDir() && f.Name()[0] != '.' {
			activatePth := filepath.Join(env.DbDir, f.Name(), "activate")
			hasActivate, err := IsFile(activatePth)
			if err != nil {
				return nil, fmt.Errorf("'%v': %v", activatePth, err)
			}
			if hasActivate {
				names = append(names, f.Name())
			}
		}
	}
	return
}

func (env *GoEnv) ActivateCode(name string) (string, error) {
	return fmt.Sprintf("source %q\ns=$?\n[ $s -ne 0 ]",
		filepath.Join(env.DbDir, name, "activate")), nil
}

func (env *GoEnv) Rm(name string, delete bool) (string, error) {
	var (
		exists bool
		err error
	)
	pth := filepath.Join(env.DbDir, name)
	exists, err = IsDir(pth)
	if !exists {
		return "", fmt.Errorf("%v: %v", pth, os.ErrNotExist)
	}

	exists, err = IsFile(pth, "activate")
	if err != nil {
		return "", fmt.Errorf("%v: %v", filepath.Join(pth, "activate"), err)
	}
	if !exists {
		return "", fmt.Errorf("%v is not GoLang enviroment.", pth)
	}

	if delete {
		return "", os.RemoveAll(pth)
	}

	trashDir := filepath.Join(env.DbDir, ".trash")
	exists, err = IsDir(trashDir)

	if err != nil {
		return "", err
	}

	if !exists {
		if err = MkdirAll(trashDir); err != nil {
			return "", fmt.Errorf("%v: %v", trashDir, err)
		}
	}

	newPth := filepath.Join(trashDir, name+"_"+TimeString(time.Now()))
	err = os.Rename(pth, newPth)

	if err != nil {
		return "", err
	}
	return newPth, nil
}
