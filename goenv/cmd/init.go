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

var initCmd = &cobra.Command{
	Use:   "init NAME...",
	Short: "Init new virtual enviroment.",
	Long: `Init new virtual enviroment on current database.

Examples:
  $ goenv init env1 env2
  $ goenv -d ~/my-goenv init env3 env4
`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return nil
		}
		env, err := goenv.NewGoEnvCmd(db, false)
		if err != nil {
			return err
		}

		return env.Init(args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
