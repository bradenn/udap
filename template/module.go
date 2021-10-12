package template

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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

func LoadEnvironment(env string, environment interface{}) error {
	err := json.Unmarshal([]byte(env), &environment)
	if err != nil {
		return err
	}

	of := reflect.ValueOf(environment)
	for i := 0; i < of.NumField(); i++ {
		err := os.Setenv(of.Field(i).String(), of.Field(i).Interface().(string))
		if err != nil {
			return err
		}
		fmt.Println(of.Field(i).String(), ": ", of.Field(i).Interface().(string))
	}
	return nil
}
