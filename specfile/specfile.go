package specfile

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const SUFFIX = ".yml"

type HostSpec struct {
	Hostname string `json:"hostname"`
	File     string `json:"file"`
	Port     int    `json:"port"`
}

type KeySpec struct {
	Path string `json:"path"`
}

type SpecData struct {
	Hosts map[string]HostSpec `json:"hosts"`
	Keys  map[string]KeySpec  `json:"keys"`
}

type SpecTemplateConfig struct {
	WithComments bool
	ExcludeKeys  bool
}

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

func ReadSpecFile(filename string) (*SpecData, error) {
	d, f := filepath.Split(filename)
	if d == "" {
		d = "."
	}
	viper.AddConfigPath(d)
	viper.SetConfigName(strings.TrimSuffix(f, SUFFIX))
	viper.SetConfigType("yaml")

	var err error
	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("Unable to read %s: %v", filename, err)
	}

	data := &SpecData{}
	err = viper.Unmarshal(data)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal config to SpecData instance: %v", err)
	}
	return data, nil
}
