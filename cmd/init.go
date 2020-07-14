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
	"os"
	"strings"

	"github.com/drognisep/sshtail/specfile"
	"github.com/spf13/cobra"
)

var withComments bool
var excludeKeys bool
var overwrite bool

const suffix string = ".yml"

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Args:  cobra.ExactArgs(1),
	Short: "Initializes a spec template with the given file name and the .yml suffix added",
	Long: `This will create a spec file showing what hosts to connect to, what file to
tail, and what keys to use (keys are optional to promote portability).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := strings.TrimSuffix(args[0], suffix) + suffix
		fmt.Printf("Creating template spec file '%s'\n", filename)
		config := &specfile.SpecTemplateConfig{WithComments: withComments, ExcludeKeys: excludeKeys}
		text, err := specfile.NewSpecTemplate(config)
		if err != nil {
			return err
		}
		if overwrite == false {
			_, err = os.Stat(filename)
			if err == nil {
				var response string
				fmt.Printf("The file already exists, do you want to replace it (Y/N)? ")
				fmt.Scanf("%s\n", &response)
				response = strings.TrimSpace(strings.ToLower(response))
				switch response {
				case "yes":
				case "y":
				default:
					fmt.Println("Canceling init operation")
					return nil
				}
			}
		}
		err = ioutil.WriteFile(filename, []byte(text), 0644)
		if err != nil {
			return fmt.Errorf("Unable to write to file %s: %v", filename, err)
		}
		fmt.Println("Spec written to file")
		return nil
	},
}

func init() {
	specCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&withComments, "with-comments", "", false, "Include comments in the template. This can be useful for understanding the format")
	initCmd.Flags().BoolVarP(&excludeKeys, "exclude-keys", "", false, "Exclude the keys section to create a portable spec file")
	initCmd.Flags().BoolVarP(&overwrite, "overwrite", "", false, "Do not check for the existence of the target file, overwrite it.")
}
