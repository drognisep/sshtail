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
package specfile

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// HostSpec identifies the hostname and port to connect to, as well as the file to tail.
type HostSpec struct {
	Hostname string `json:"hostname" yaml:"hostname"`
	File     string `json:"file" yaml:"file"`
	Port     int    `json:"port" yaml:"port"`
}

// KeySpec specifies the path to the SSH key to be used for the host named by the SpecData.Keys map key.
type KeySpec struct {
	Path string `json:"path" yaml:"path"`
}

// SpecData encapsulates runtime parameters for SSH tailing.
type SpecData struct {
	Hosts map[string]HostSpec `json:"hosts" yaml:"hosts"`
	Keys  map[string]KeySpec  `json:"keys" yaml:"keys"`
}

// SpecTemplateConfig config
type SpecTemplateConfig struct {
	WithComments bool
	ExcludeKeys  bool
}

// NewSpecTemplate creates a new spec template with the given configuration parameters.
func NewSpecTemplate(config *SpecTemplateConfig) (string, error) {
	templateString := `{{- if .WithComments}}# Hosts and files to tail
{{end}}hosts:
  host1:
    hostname: remote-host-1
    file: /var/log/syslog
    {{if .WithComments}}# Default SSH port
    {{end}}port: 22
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

// ReadSpecFile attempts to read SpecData from the specified file.
func ReadSpecFile(filename string) (*SpecData, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	specData := &SpecData{}
	err = yaml.Unmarshal(data, specData)
	if err != nil {
		return nil, err
	}

	return specData, nil
}
