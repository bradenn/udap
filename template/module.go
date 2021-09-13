package template

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
}

type Function func(string) (string, error)

type Config struct {
	path string
}

func (c Config) realPath() string {
	filename := fmt.Sprintf("./plugins/%s/config.json", c.path)
	return filename
}

func NewConfig(path string) Config {
	config := Config{
		path: path,
	}

	return config
}

func (c *Config) writeFile(object map[string]string) (err error) {

	marshal, err := json.Marshal(object)
	if err != nil {
		return err
	}

	jsonFile, err := os.OpenFile(c.realPath(), os.O_RDWR|os.O_CREATE, 0644)
	defer jsonFile.Close()
	if err != nil {
		return err
	}
	err = jsonFile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = jsonFile.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(jsonFile, string(marshal))
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) readFile() (object map[string]string, err error) {
	object = map[string]string{}

	jsonFile, err := os.OpenFile(c.realPath(), os.O_RDONLY|os.O_CREATE, 0644)
	defer jsonFile.Close()
	if err != nil {
		return object, err
	}

	all, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(all, &object)
	if err != nil {
		return object, err
	}

	return object, nil
}

func (c *Config) Set(key string, model interface{}) {
	file, err := c.readFile()
	if err != nil {
		file = map[string]string{}
	}

	marshal, err := json.Marshal(model)
	if err != nil {
		return
	}

	file[key] = string(marshal)

	err = c.writeFile(file)
	if err != nil {
		return
	}

}

func (c *Config) IsSet(key string) bool {
	file, err := c.readFile()
	if err != nil {
		return false
	}

	if file[key] == "" {
		return false
	}

	return true
}

func (c *Config) Get(key string, model interface{}) error {
	file, err := c.readFile()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = json.Unmarshal([]byte(file[key]), &model)
	if err != nil {
		return err
	}

	return nil
}

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
