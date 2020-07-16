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
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os/user"
	"path"

	"gopkg.in/yaml.v3"
)

const DEFAULT_SSH_PORT int = 22

// HostSpec identifies the hostname and port to connect to, as well as the file to tail.
type HostSpec struct {
	Hostname string `json:"hostname" yaml:"hostname"`
	Username string `json:"username" yaml:"username"`
	File     string `json:"file" yaml:"file"`
	Port     int    `json:"port" yaml:"port"`
}

// Validate checks the HostSpec for errors and sets reasonable defaults.
func (h *HostSpec) Validate() error {
	if h.Hostname == "" {
		return errors.New("Host spec cannot have a blank hostname")
	}
	if h.Username == "" {
		u, err := user.Current()
		if err != nil {
			// May need to check for sudo user on linux. Not going to support
			// edge cases like this initially.
			return errors.New("Unable to determine current user")
		}
		h.Username = u.Username
	}
	if h.File == "" {
		return errors.New("Host spec cannot have a blank file")
	}
	if h.Port == 0 {
		h.Port = DEFAULT_SSH_PORT
	}
	return nil
}

// KeySpec specifies the path to the SSH key to be used for the host named by the SpecData.Keys map key.
type KeySpec struct {
	Path string `json:"path" yaml:"path"`
}

// Validate checks the KeySpec for errors and sets reasonable defaults.
func (k *KeySpec) Validate() error {
	if k.Path == "" {
		k.Path = DefaultSSHKeyPath()
	}
	return nil
}

// SpecData encapsulates runtime parameters for SSH tailing.
type SpecData struct {
	Hosts map[string]HostSpec `json:"hosts" yaml:"hosts"`
	Keys  map[string]KeySpec  `json:"keys" yaml:"keys"`
}

// Validate checks the SpecData for errors and sets reasonable defaults.
func (s *SpecData) Validate() error {
	if s.Hosts == nil || len(s.Hosts) == 0 {
		return errors.New("Host spec must have at least one definition")
	}
	if s.Keys != nil && len(s.Keys) > 0 {
		for k, v := range s.Keys {
			err := v.Validate()
			if err != nil {
				return fmt.Errorf("Key spec %s: %v", k, err)
			}
		}
	} else {
		s.Keys = map[string]KeySpec{}
	}

	for k, v := range s.Hosts {
		err := v.Validate()
		if err != nil {
			return fmt.Errorf("Host spec %s: %v", k, err)
		}
		_, found := s.Keys[k]
		if !found {
			s.Keys[k] = KeySpec{DefaultSSHKeyPath()}
		}
	}

	hostsLen := len(s.Hosts)
	keysLen := 0
	if s.Keys != nil {
		keysLen = len(s.Keys)
	}

	if keysLen != 0 && hostsLen != keysLen {
		fmt.Println("Warning: The number of host entries does not match the number of keys entries")
	}

	return nil
}

// ConfigFileData is the format for the home config file.
type ConfigFileData struct {
	DefaultKey KeySpec `json:"defaultKey" yaml:"defaultKey"`
}

func defaultSSHKeyPath() string {
	u, _ := user.Current()
	return path.Join(u.HomeDir, ".ssh", "id_rsa")
}

// DefaultSSHKeyPath returns the config specified value or the default path of ~/.ssh/id_rsa
func DefaultSSHKeyPath() string {
	var ks KeySpec
	c, err := ConfigFile()
	if err != nil || c == nil || c.DefaultKey == ks {
		ks = KeySpec{defaultSSHKeyPath()}
	} else {
		ks = c.DefaultKey
	}
	return ks.Path
}

// ConfigFile reads the default config from the active user's home directory.
func ConfigFile() (*ConfigFileData, error) {
	u, _ := user.Current()
	p := path.Join(u.HomeDir, ".sshtail.yaml")
	confData, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file: %v", err)
	}
	confFileData := &ConfigFileData{}
	err = yaml.Unmarshal(confData, confFileData)
	if err != nil {
		return nil, fmt.Errorf("Config file is not a valid format: %v", err)
	}
	return confFileData, nil
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
    {{if .WithComments}}# Excluding the username here will default it to the current user name
    {{end}}file: /var/log/syslog
    {{if .WithComments}}# Default SSH port
    {{end}}port: 22
  host2:
    hostname: remote-host-2
    username: me
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
