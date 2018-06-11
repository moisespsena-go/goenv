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
	"github.com/spf13/cobra"
	"github.com/moisespsena/go-goenv"
	"os"
)

// rmCmd represents the rm command
var backupCmd = &cobra.Command{
	Use:   "backup NAME [DEST]",
	Short: "Create backup .tar.gz file for the virtualenv with have NAME.",
	Long: `Create backup .tar.gz file for the virtualenv with have NAME.
Examples:
  $ goenv backup teste target.tar.gz

  backup to stdout:

  $ goenv backup teste -
`,
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

		if err != nil {
			return err
		}

		options := &goenv.BackupOptions{}

		if len(args) == 1 {
			options.DefaultBackup = true
		} else if args[1] == "-" {
			options.Writer = os.Stdout
		} else {
			options.Target = args[1]
		}

		return env.Backup(args[0], options)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
