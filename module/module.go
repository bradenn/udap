package module

import (
	"fmt"
	"os"
	"plugin"
	"strings"
)

type AModule struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
}

var modules map[string]Module

type Error struct {
	error string
}

func NewError(message string) Error {
	return Error{error: message}
}

func (e Error) Error() string {
	return e.error
}

func List() (m []Module, err error) {
	for _, module := range modules {
		m = append(m, module)
	}
	return m, err
}

func Get(identifier string) (Module, error) {
	if modules[identifier] == nil {
		return nil, NewError(fmt.Sprintf("Could not find module '%s'", identifier))
	}
	return modules[identifier], nil
}

type Module interface {
	Name() string
	Description() string
	Functions() []string
	Run(name string, payload ...interface{}) interface{}
	Get()
	Pub()
}

func init() {
	var err error
	modules, err = loadModules("./plugins")
	if err != nil {
		panic(err)
	}
}

func loadModules(path string) (modules map[string]Module, err error) {
	modules = map[string]Module{}
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range dir {
		if entry.IsDir() {
			filename := fmt.Sprintf("./plugins/%s/%s.so", entry.Name(), entry.Name())
			module, err := loadModule(filename)
			if err != nil {
				return modules, err
			}

			modules[strings.ToLower(module.Name())] = module
			fmt.Printf("Loaded module '%s' with %d functions.\n", module.Name(), len(module.Functions()))
		}
	}

	return modules, err
}

func loadModule(path string) (Module, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	lookup, err := p.Lookup("Module")
	if err != nil {
		return nil, err
	}

	return lookup.(Module), nil
}
