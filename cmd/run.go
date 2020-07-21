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
	"fmt"
	"os"

	"github.com/drognisep/sshtail/specfile"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Args:  cobra.ExactArgs(1),
	Short: "Runs a spec file to connect to multiple hosts and tail the files specified",
	Long: `Spec files have the extension .spec. A template can be created with
	sshtail spec init your-spec-name-here`,
	RunE: func(cmd *cobra.Command, args []string) error {
		specData, err := specfile.ReadSpecFile(args[0])
		if err != nil {
			return fmt.Errorf("Unable to parse config file '%s': %v", args[0], err)
		}
		writer, err := specfile.NewConsolidatedWriter(specData, os.Stdout)
		if err != nil {
			return err
		}
		writer.Start()
		return nil
	},
}

func init() {
	specCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
