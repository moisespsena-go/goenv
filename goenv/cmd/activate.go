// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"os"
	"fmt"
	"path/filepath"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate NAME",
	Short: "Activate the virtualenv with NAME.",
	Long: `Activate the virtualenv with NAME.
Examples:
  $ eval $(goenv activate teste)
  $ eval $(goenv -d ~/my-goenv activate teste)
`,
	Args:cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ok, err := IsDir(db, args[0])
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintf(os.Stderr, "'%v': Database isn't initialized.\n", db)
			return nil
		}
		defer func() {
			os.Stdout.Sync()
			os.Stderr.Sync()
		}()

		fmt.Fprintf(os.Stdout, "source %q\ns=$?\n[ $s -ne 0 ]",
			filepath.Join(db, args[0], "activate"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
