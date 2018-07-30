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
	"github.com/moisespsena/go-goenv"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Generate sources for custom prompt commands.",
	Long: `Generate sources for custom prompt commands.
Examples:
  $ goenv setup | tee -a ~/.bashrc
  $ source ~/.bashrc
    or
  $ eval $(goenv setup)
  
Commands available:
  - goenv-init
    Example (see for 'init' sub command):
        $ goenv-init env1

  - goenv-activate
    Example:
        $ goenv-activate env1

  - goenv-die
    Clean all goenv commands from 'goenv setup'

    Example:
        $ goenv-die
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := goenv.NewGoEnvCmd(db, false)
		if err != nil {
			return err
		}
		return env.Setup()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
