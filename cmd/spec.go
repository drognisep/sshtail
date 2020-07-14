/*
Copyright © 2020 Joseph Saylor <doug@saylorsolutions.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// specCmd represents the spec command
var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Root command for operations involving spec files",
	Long: `Spec files (*.spec) are used to store information about multiple hosts, what
files to tail, and optionally what keys to use to connect with the server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Must specify additional commands")
	},
}

func init() {
	rootCmd.AddCommand(specCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// specCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// specCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
