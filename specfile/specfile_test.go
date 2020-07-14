package specfile

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

var hostAndKeysData SpecData = SpecData{
	map[string]HostSpec{
		"host1": HostSpec{"remote-host-1", "/var/log/syslog", 22},
		"host2": HostSpec{"remote-host-2", "/var/log/syslog", 22},
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
		"host1": HostSpec{"remote-host-1", "/var/log/syslog", 22},
		"host2": HostSpec{"remote-host-2", "/var/log/syslog", 22},
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
    file: /var/log/syslog
    # Default SSH port
    port: 22
  host2:
    hostname: remote-host-2
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
		"host1": HostSpec{"remote-host-1", "/var/log/syslog", 22},
		"host2": HostSpec{"remote-host-2", "/var/log/syslog", 22},
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
