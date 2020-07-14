package data

type HostSpec struct {
	Hostname string `json:"hostname"`
	File     string `json:"file"`
	Port     int    `json:"port"`
}

type KeySpec struct {
	Path string `json:"path"`
}

type SpecFile struct {
	Hosts map[string]HostSpec `json:"hosts"`
	Keys  map[string]KeySpec  `json:"keys"`
}
