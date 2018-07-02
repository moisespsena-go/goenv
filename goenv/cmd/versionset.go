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
	"github.com/moisespsena/go-error-wrap"
)

// activateCmd represents the activate command
var versionsSetCmd = &cobra.Command{
	Use:   "set VERSION ENV_NAME...",
	Short: "Set version to env",
	Args:cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := goenv.NewGoEnv(db, false)
		if err != nil {
			return errwrap.Wrap(err, "New Env")
		}
		vs := goenv.NewGoVersions(env)
		versionName, args := args[0], args[1:]
		for _, envName := range args {
			err = vs.Set(versionName, envName)
			if err != nil {
				return errwrap.Wrap(err, "Set version to %q", envName)
			}
		}
		return nil
	},
}

func init() {
	versionsCmd.AddCommand(versionsSetCmd)
}
