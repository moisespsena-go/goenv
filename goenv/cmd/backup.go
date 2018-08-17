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
	"os"

	"github.com/moisespsena/go-error-wrap"
	"github.com/moisespsena/go-goenv"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup NAME [DEST]",
	Short: "Create backup .tar.gz file for the virtualenv with have NAME.",
	Long: `Create backup .tar.gz file for the virtualenv with have NAME.
Examples:
  $ goenv backup teste target.tar.gz

  backup to stdout:

  $ goenv backup teste -

Exclude patterns from backup:
  With Args:

  $ goenv backup teste -e ".git" -e "node_modules" -e "*.swp"

  Global:

  $ ENV_PATH=$(goenv path teste)
  $ mkdir $ENV_PATH/.goenv_settings
  $ echo ".git\nnode_modules\n*.swp" > $ENV_PATH/.goenv_settings/backup_exclude
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

		options := &goenv.BackupOptions{}
		exclude, err := cmd.PersistentFlags().GetStringSlice("exclude")
		if err != nil {
			return errwrap.Wrap(err, "Flag EXCLUDE")
		}

		if len(exclude) > 0 {
			err = options.Patterns.Append(exclude...)
			if err != nil {
				return errwrap.Wrap(err, "Exclude patterns.")
			}
		}

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
	backupCmd.PersistentFlags().StringSliceP("exclude", "e", nil,
		"Excludes using GLOB. See https://github.com/gobwas/glob for patthern help.")
	rootCmd.AddCommand(backupCmd)
}
