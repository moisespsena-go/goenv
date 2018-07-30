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

	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Returns the current database path.",
	Long: `Returns the current database path.

For change current database path:

With command paramenter
-----------------------

For use custom database path, set the '-d' flag:

	goenv -d ~/custom-db args...

Examples:

	goenv -d ~/custom-db init env2
	goenv -d ~/custom-db ls

With enviroment variable
-----------------------

Set the enviroment variable:
 
	export GOENVDB=~/custom-db

or

	GOENVDB=~/custom-db goenv args...
`,
	Run: func(cmd *cobra.Command, args []string) {
		os.Stdout.WriteString(db + "\n")
		os.Stdout.Sync()
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
}
