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

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore NAME BACKUP.tar.gz",
	Short: "Restore backup.tar.gz file to the virtualenv with have NAME.",
	Long: `Restore backup.tar.gz file to the virtualenv with have NAME.
Examples:
  Restore from file.
  	$ goenv restore env1.tar.gz

  Restore from stdout:
  	$ cat backup.tar.gz | goenv restore

  Restore with another name:
	$ goenv restore -n env2 backup.tar.gz
  	$ cat env1.tar.gz | goenvrestore -n env2

`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := goenv.NewGoEnvCmd(db, true)
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}

		options := &goenv.RestoreOptions{}

		if len(args) == 1 {
			options.Source = args[0]
		} else {
			fi, err := os.Stdin.Stat()
			if err != nil {
				panic(err)
			}
			if fi.Mode()&os.ModeNamedPipe == 0 {
				options.Reader = os.Stdin
			}
		}

		options.Name, err = cmd.PersistentFlags().GetString("name")
		if err != nil {
			return err
		}
		options.Update, err = cmd.PersistentFlags().GetBool("update")
		if err != nil {
			return err
		}
		options.OverWrite, err = cmd.PersistentFlags().GetBool("overwrite")
		if err != nil {
			return err
		}
		options.Archive, err = cmd.PersistentFlags().GetBool("archive")
		if err != nil {
			return err
		}
		options.Verbose, err = cmd.PersistentFlags().GetBool("verbose")
		if err != nil {
			return err
		}
		options.Trial, err = cmd.PersistentFlags().GetBool("dry-run")
		if err != nil {
			return err
		}

		return env.Restore(options)
	},
}

func init() {
	restoreCmd.PersistentFlags().BoolP("overwrite", "o", false,
		"Overwrite enviroment if exists (default is false).")
	restoreCmd.PersistentFlags().BoolP("update", "u", false,
		"Update enviroment if exists (default is false).")
	restoreCmd.PersistentFlags().BoolP("archive", "a", false,
		"Archive only. Not compressed with GZip.")
	restoreCmd.PersistentFlags().BoolP("verbose", "v", false,
		"Print names while restoring.")
	restoreCmd.PersistentFlags().BoolP("dry-run", "D", false,
		"Perform a trial run with no changes made.")
	restoreCmd.PersistentFlags().StringP("name", "n", "",
		"Name after restored.")
	rootCmd.AddCommand(restoreCmd)
}
