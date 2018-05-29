package go_goenv

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
		fmt.Fprintf(os.Stdout, "Activate it using:\n  $ source '%v'\n    or\n  $ eval $(%v activate %v)\n\n",
			filepath.Join(pth, "activate"), os.Args[0], name)
	}

	fmt.Fprintln(os.Stdout, "done")
	os.Stdout.Sync()
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
