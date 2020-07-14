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
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

var hostAndKeysData SpecData = SpecData{
	map[string]HostSpec{
		"host1": HostSpec{"remote-host-1", "", "/var/log/syslog", 22},
		"host2": HostSpec{"remote-host-2", "me", "/var/log/syslog", 22},
	},
	map[string]KeySpec{
		"host1": KeySpec{"~/.ssh/id_rsa"},
		"host2": KeySpec{"~/.ssh/id_rsa"},
	},
}

const defaultSpecText string = `hosts:
  host1:
    hostname: remote-host-1
    file: /var/log/syslog
    port: 22
  host2:
    hostname: remote-host-2
    username: me
    file: /var/log/syslog
    port: 22
keys:
  host1:
    path: ~/.ssh/id_rsa
  host2:
    path: ~/.ssh/id_rsa
`

var commentHostAndKeysData SpecData = SpecData{
	map[string]HostSpec{
		"host1": HostSpec{"remote-host-1", "", "/var/log/syslog", 22},
		"host2": HostSpec{"remote-host-2", "me", "/var/log/syslog", 22},
	},
	map[string]KeySpec{
		"host1": KeySpec{"~/.ssh/id_rsa"},
		"host2": KeySpec{"~/.ssh/id_rsa"},
	},
}

const commentSpecText string = `# Hosts and files to tail
hosts:
  host1:
    hostname: remote-host-1
    # Excluding the username here will default it to the current user name
    file: /var/log/syslog
    # Default SSH port
    port: 22
  host2:
    hostname: remote-host-2
    username: me
    file: /var/log/syslog
    port: 22
# This section is optional for portability
keys:
  host1:
    # Defaults to this value
    path: ~/.ssh/id_rsa
  host2:
    # If all of these values are the same, then 'sshtail usekey' may be more convenient.
    path: ~/.ssh/id_rsa
`

var hostData SpecData = SpecData{
	map[string]HostSpec{
		"host1": HostSpec{"remote-host-1", "", "/var/log/syslog", 22},
		"host2": HostSpec{"remote-host-2", "me", "/var/log/syslog", 22},
	},
	nil,
}

const noKeysSpecText string = `hosts:
  host1:
    hostname: remote-host-1
    file: /var/log/syslog
    port: 22
  host2:
    hostname: remote-host-2
    username: me
    file: /var/log/syslog
    port: 22
`

func TestInitDefaultSpec(t *testing.T) {
	want := defaultSpecText
	got, err := NewSpecTemplate(&SpecTemplateConfig{false, false})
	if err != nil {
		t.Errorf("Failed to create spec template: %v", err)
	}
	if got != defaultSpecText {
		t.Errorf("Got:\n%v\nWanted:\n%v", got, want)
	}
}

func TestInitWithComments(t *testing.T) {
	want := commentSpecText
	got, err := NewSpecTemplate(&SpecTemplateConfig{true, false})

	if err != nil {
		t.Errorf("Failed to create spec template: %v", err)
	}

	if got != want {
		t.Errorf("Got:\n%v\nWanted:\n%v", got, want)
	}
}

func TestInitNoKeys(t *testing.T) {
	want := noKeysSpecText
	got, err := NewSpecTemplate(&SpecTemplateConfig{false, true})

	if err != nil {
		t.Errorf("Failed to create spec template: %v", err)
	}

	if got != want {
		t.Errorf("Got:\n%v\nWanted:\n%v", got, want)
	}
}

func TestReadDefault(t *testing.T) {
	constVal, err := json.Marshal(hostAndKeysData)
	if err != nil {
		t.Error("Unable to serialize testing value")
	}

	text, err := NewSpecTemplate(&SpecTemplateConfig{false, false})
	if err != nil {
		t.Errorf("Failed to create spec template: %v", err)
	}

	ioutil.WriteFile("testDefault.yml", []byte(text), 0644)
	defer os.Remove("testDefault.yml")
	data, err := ReadSpecFile("testDefault.yml")
	if err != nil {
		t.Errorf("Unable to read from file: %v", err)
	}

	actVal, err := json.Marshal(data)
	if err != nil {
		t.Error("Unable to serialize actual value")
	}
	got, want := string(actVal), string(constVal)

	if got != want {
		t.Errorf("Got:\n%s\nWant:\n%s\n", got, want)
	}
}

func TestReadCommented(t *testing.T) {
	constVal, err := json.Marshal(commentHostAndKeysData)
	if err != nil {
		t.Error("Unable to serialize testing value")
	}

	text, err := NewSpecTemplate(&SpecTemplateConfig{true, false})
	if err != nil {
		t.Errorf("Failed to create spec template: %v", err)
	}

	ioutil.WriteFile("testCommented.yml", []byte(text), 0644)
	defer os.Remove("testCommented.yml")
	data, err := ReadSpecFile("testCommented.yml")
	if err != nil {
		t.Errorf("Unable to read from file: %v", err)
	}

	actVal, err := json.Marshal(data)
	if err != nil {
		t.Error("Unable to serialize actual value")
	}
	got, want := string(actVal), string(constVal)

	if got != want {
		t.Errorf("Got:\n%s\nWant:\n%s\n", got, want)
	}
}

func TestReadHostOnly(t *testing.T) {
	constVal, err := json.Marshal(hostData)
	if err != nil {
		t.Error("Unable to serialize testing value")
	}

	text, err := NewSpecTemplate(&SpecTemplateConfig{false, true})
	if err != nil {
		t.Errorf("Failed to create spec template: %v", err)
	}

	ioutil.WriteFile("testHostOnly.yml", []byte(text), 0644)
	defer os.Remove("testHostOnly.yml")
	data, err := ReadSpecFile("testHostOnly.yml")
	if err != nil {
		t.Errorf("Unable to read from file: %v", err)
	}

	actVal, err := json.Marshal(data)
	if err != nil {
		t.Error("Unable to serialize actual value")
	}
	got, want := string(actVal), string(constVal)

	if got != want {
		t.Errorf("Got:\n%s\nWant:\n%s\n", got, want)
	}
}

func TestValidateHost(t *testing.T) {
	errorList := []HostSpec{
		HostSpec{"", "me", "file", 22}, // No host
		HostSpec{"host", "me", "", 22}, // No file
	}

	for i, h := range errorList {
		err := h.Validate()
		if err == nil {
			t.Errorf("'errorList[%d]' should not have passed validation", i)
		}
	}
}

func TestValueDefaultHost(t *testing.T) {
	missingUser := HostSpec{"host", "", "file", 22}
	missingPort := HostSpec{"host", "me", "file", 0}
	var err error

	err = missingUser.Validate()
	if err != nil {
		t.Errorf("Should not have received error for missing username: %v", err)
	}
	if missingUser.Username == "" {
		t.Error("Username was not defaulted")
	}

	err = missingPort.Validate()
	if err != nil {
		t.Errorf("Should not have received error for missing port: %v", err)
	}
	if missingPort.Port != DEFAULT_SSH_PORT {
		t.Error("Default port was not set")
	}
}
