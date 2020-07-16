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
	"io/ioutil"
	"os/user"
	"path"

	"github.com/drognisep/sshtail/specfile"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

// usekeyCmd represents the usekey command
var usekeyCmd = &cobra.Command{
	Use:   "usekey",
	Args:  cobra.ExactArgs(1),
	Short: "Specifies the default SSH key",
	Long: `This is useful for situations where the same SSH key is used to connect to
multiple nodes in a cluster.

Instead of specifying the same key multiple times, the default key can be used
for all of them. The default key path will be saved to your config file for
later use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, _ := specfile.ConfigFile()
		if config == nil {
			// Failed to read in the file. Could be not created yet.
			config = &specfile.ConfigFileData{}
		}

		config.DefaultKey = specfile.KeySpec{Path: args[0]}

		newConfig, err := yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("Internal error: %v", err)
		}
		u, _ := user.Current()
		err = ioutil.WriteFile(path.Join(u.HomeDir, ".sshtail.yaml"), newConfig, 0644)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(usekeyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usekeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usekeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
