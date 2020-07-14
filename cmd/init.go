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
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var withComments bool
var excludeKeys bool
var overwrite bool

const SUFFIX string = ".yml"

type templateConfig struct {
	WithComments bool
	ExcludeKeys  bool
}

func specTemplate(config *templateConfig) (string, error) {
	templateString := `{{- if .WithComments}}# Hosts and files to tail
{{end}}hosts:
  host1:
    hostname: remote-host-1
	file: /var/log/syslog
	{{if .WithComments}}# Default SSH port{{end}}
	port: 22
  host2:
    hostname: remote-host-2
	file: /var/log/syslog
	port: 22
{{if not .ExcludeKeys}}{{if .WithComments}}# This section is optional for portability
{{end}}keys:
  host1:
    {{if .WithComments}}# Defaults to this value
    {{end}}path: ~/.ssh/id_rsa
  host2:
    {{if .WithComments}}# If all of these values are the same, then 'sshtail usekey' may be more convenient.
    {{end}}path: ~/.ssh/id_rsa
{{end}}`
	t, err := template.New("spec-template").Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("Unable to parse template: %v", err)
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, config); err != nil {
		return "", fmt.Errorf("Unable to generate template file contents: %v", err)
	}

	return buf.String(), nil
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Args:  cobra.ExactArgs(1),
	Short: "Initializes a spec template (YAML) with the given file name",
	Long: `This will create a spec file showing what hosts to connect to, what file to
tail, and what keys to use (keys are optional to promote portability).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := strings.TrimSuffix(args[0], SUFFIX) + SUFFIX
		fmt.Printf("Creating template spec file '%s'\n", filename)
		config := &templateConfig{withComments, excludeKeys}
		text, err := specTemplate(config)
		if err != nil {
			return err
		}
		// fmt.Printf("Generated template:\n%s", text)
		if overwrite == false {
			_, err = os.Stat(filename)
			if err == nil {
				var response string
				fmt.Printf("The file already exists, do you want to replace it (Y/N)? ")
				fmt.Scanf("%s\n", &response)
				response = strings.ToLower(response)
				switch response {
				case "yes":
				case "y":
				default:
					return errors.New("Canceling init operation")
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().BoolVarP(&withComments, "with-comments", "", false, "Include comments in the template. This can be useful for understanding the format")
	initCmd.Flags().BoolVarP(&excludeKeys, "exclude-keys", "", false, "Exclude the keys section to create a portable spec file")
	initCmd.Flags().BoolVarP(&overwrite, "overwrite", "", false, "Do not check for the existence of the target file, overwrite it.")
}
