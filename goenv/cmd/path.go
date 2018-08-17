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

package cmd

import (
	"github.com/moisespsena/go-error-wrap"
	"github.com/moisespsena/go-goenv"
	"github.com/spf13/cobra"
	"os"
)

var pathCmd = &cobra.Command{
	Use:   "path NAME...",
	Short: "Print env path",
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.MinimumNArgs(1)(cmd, args)
		if err != nil {
			return err
		}
		return cobra.MaximumNArgs(2)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := goenv.NewGoEnvCmd(db, true)
		if err != nil {
			return err
		}

		var pth string

		for i, name := range args {
			pth, err = env.Env.GetCheck(name)
			if err != nil {
				return errwrap.Wrap(err, "Arg %d: %q", i, name)
			}

			if _, err = os.Stdout.WriteString(pth + "\n"); err != nil {
				return err
			}
			if err = os.Stdout.Sync(); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)
}
