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

	"github.com/moisespsena/go-error-wrap"
	"github.com/moisespsena/go-goenv"
	"github.com/spf13/cobra"
)

var versionsAvailableCmd = &cobra.Command{
	Use:   "available [TERM...]",
	Short: "List all available GoLang versions",
	Long: `List all available GoLang versions
The TERM is Glob (https://github.com/gobwas/glob) expression.

Examples:
  $ goenv available
  $ goenv available 1.1*
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := goenv.NewGoEnv(db, false)
		if err != nil {
			return errwrap.Wrap(err, "New Env")
		}
		v := goenv.NewGoVersions(env)
		items, err := v.Available(args...)
		if err != nil {
			return err
		}
		fmt.Println(pad("Name"), pad("URL", 55), "Root")

		for _, v := range items {
			fmt.Println(pad(v.Name), pad(v.DownloadUrl(), 55), v.Root)
		}
		return nil
	},
}

func init() {
	versionsCmd.AddCommand(versionsAvailableCmd)
}
