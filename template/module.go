package template

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
}

type Mod interface {
	// InitInstance is called when a new instance is created, the string returned is saved to the instance env.
	InitInstance() (string, error)
	// Initialize is run when the module is loaded, it receives global environment variables.
	Initialize(env string)
	// Metadata returns basic information about the module.
	Metadata() Metadata
	// Poll is called when data is requested, env is the environment saved from Init
	Poll(v string) (string, error)
	// Run will request a function to be called, the outputs are returned.
	Run(v string, action string) (string, error)
}

type Function func(string) (string, error)

type Module struct {
	metadata   Metadata
	functions  map[string]Function
	config     Config
	onEnable   func()
	instanceId string
}

func NewModule(metadata Metadata, functions map[string]Function, onEnable func()) Module {
	configString := strings.ToLower(metadata.Name)
	config := NewConfig(configString)
	return Module{metadata: metadata, functions: functions, onEnable: onEnable, config: config}
}

func (m *Module) Metadata() Metadata {
	return m.metadata
}

func (m *Module) GetInstance() string {
	return m.instanceId
}

func (m *Module) GetConfig() Config {
	return m.config
}

func (m *Module) Run(s string) (string, error) {
	if m.functions[s] == nil {
		return "", fmt.Errorf("function does not exist")
	}
	return m.functions[s](s)
}

func (m *Module) Configure(data []byte, instanceId string) {
	m.instanceId = instanceId
	raw := map[string]string{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return
	}
	for s, message := range raw {
		err := os.Setenv(s, message)
		if err != nil {
			return
		}
	}
	m.onEnable()
}

func (m *Module) Functions() (functions []string) {
	for s := range m.functions {
		functions = append(functions, s)
	}
	return functions
}
