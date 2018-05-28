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
	"fmt"

	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all virtual enviroments on current database.",
	Long: `List all virtual enviroments on current database.

Examples:
  $ goenv ls
  env1
  env2

  $ goenv -d ~/my-goenv ls
  env3
  env4
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ok, err := IsDir(db)
		if err != nil {
			return err
		}
		if !ok {
			fmt.Fprintf(os.Stderr, "'%v': Database isn't initialized.", db)
			return nil
		}
		defer func() {
			os.Stdout.Sync()
			os.Stderr.Sync()
		}()

		files, err := ioutil.ReadDir(db)
		if err != nil {
			return fmt.Errorf("'%v': %v", db, err)
		}
		var has bool
		for _, f := range files {
			if f.IsDir() {
				has = true
				fmt.Fprintf(os.Stdout, filepath.Join(db, f.Name()))
			}
		}

		if !has {
			fmt.Fprintf(os.Stderr, "'%v': Database directory is empty.", db)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
