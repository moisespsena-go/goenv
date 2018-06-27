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
	"strings"

	"github.com/moisespsena/go-goenv"
	"github.com/spf13/cobra"
)

// activateCmd represents the activate command
var versionsAvailableCmd = &cobra.Command{
	Use:   "available",
	Short: "List all available versions of golang.",
	RunE: func(cmd *cobra.Command, args []string) error {
		//env, err := goenv.NewGoEnvCmd(db, false)
		pad := func(v string) string {
			if len(v) > 9 {
				return v
			}
			return v + strings.Repeat(" ", 10-len(v))
		}
		v := goenv.GoVersions{}
		items, err := v.Available()
		if err != nil {
			return err
		}
		for _, v := range items {
			fmt.Println(pad(v.Title), v.DownloadUrl)
		}
		return nil
	},
}

func init() {
	versionsCmd.AddCommand(versionsAvailableCmd)
}
