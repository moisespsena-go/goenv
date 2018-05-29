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
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all virtual enviroments on current database.",
	Long: `List all virtual enviroments on current database.

Examples:
  $ go-goenv ls
  env1
  env2

  $ go-goenv -d ~/my-go-goenv ls
  env3
  env4
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := go_goenv.NewGoEnvCmd(db, true)
		if err != nil {
			return err
		}

		return env.Ls()
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
