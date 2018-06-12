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
	"os"
	"path/filepath"
)

type GoEnvCmd struct {
	env *GoEnv
}

func NewGoEnvCmd(dbDir string, check bool) (envCmd *GoEnvCmd, err error) {
	env, err := NewGoEnv(dbDir, check)
	if err != nil {
		return nil, err
	}
	return &GoEnvCmd{env}, nil
}

func (cmd *GoEnvCmd) Setup() error {
	os.Stdout.WriteString(`##############################
## - BEGIN GOENV COMMANDS - ##
##############################
goenv-activate () {
 eval $(goenv activate $1)
 return $?
}
goenv-init () {
 goenv init "$@"
 return $?
}
goenv-die () {
 unset -f goenv-activate
 unset -f goenv-init
 unset -f goenv-die
}
##############################
### - END GOENV COMMANDS - ###
##############################
`)
	return nil
}

func (cmd *GoEnvCmd) Ls() error {
	names, err := cmd.env.Ls()
	if err != nil {
		return err
	}

	if len(names) == 0 {
		fmt.Fprintf(os.Stderr, "'%v': Database directory is empty.", cmd.env.DbDir)
	} else {
		for _, name := range names {
			os.Stdout.WriteString(name + "\n")
		}
	}
	return nil
}

func (cmd *GoEnvCmd) Init(names ...string) (err error) {
	defer func() {
		os.Stdout.Sync()
		os.Stderr.Sync()
	}()

	var pth string

	for _, name := range names {
		pth = filepath.Join(cmd.env.DbDir, name)
		fmt.Fprintf(os.Stdout, "Initializing virtual enviroment %q on %q...\n", name, pth)
		err = cmd.env.Init(name)
		if err != nil {
			return
		}
		fmt.Fprintf(os.Stdout, `Activate it using:
  $ goenv-activate ` + name + `
    or
  $ source '%v'
    or
  $ eval $(%v activate ` + name + `)

`,
			filepath.Join(pth, "activate"), os.Args[0])
	}

	fmt.Fprintln(os.Stdout, "done.")
	os.Stdout.Sync()
	return nil
}

func (cmd *GoEnvCmd) Update(envs []string) error {
	if len(envs) == 0 {
		return fmt.Errorf("No enviroments names informed.")
	}
	for _, envName := range envs {
		_, err := cmd.env.GetPath(envName, false)
		if err != nil {
			return err
		}
		err = cmd.env.Init(envName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cmd *GoEnvCmd) UpdateAll() error {
	envs, err := cmd.env.Ls()
	if err != nil {
		return err
	}
	for _, envName := range envs {
		err = cmd.env.Init(envName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cmd *GoEnvCmd) ActivateCode(name string) error {
	code, err := cmd.env.ActivateCode(name)
	if err != nil {
		return err
	}
	os.Stdout.WriteString(code)
	os.Stdout.Sync()
	return nil
}

func (cmd *GoEnvCmd) Rm(name string, delete bool) error {
	pth, err := cmd.env.Rm(name, delete)
	if err != nil {
		return err
	}
	if delete {
		fmt.Fprintf(os.Stdout, "GoLang Enviroment %q [%v] removed.\n", name, filepath.Join(cmd.env.DbDir, name))
		return nil
	}
	fmt.Fprintf(os.Stdout, "GoLang Enviroment %q moved from %q to %q\n", name, filepath.Join(cmd.env.DbDir, name),
		pth)
	os.Stdout.Sync()
	return nil
}
