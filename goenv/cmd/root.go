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
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var db string

var rootCmd = &cobra.Command{
	Use:   "goenv",
	Short: "The virtual enviroments manager for Go!",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if db != "" {
			db, err = homedir.Expand(db)
		}
		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	defaultDb := os.Getenv("GOENVDB")
	if defaultDb == "" {
		defaultDb = "~/.goenv"
	}
	rootCmd.PersistentFlags().StringVarP(&db, "db", "d", defaultDb, "Database directory (default is $HOME/.goenv).")
}
