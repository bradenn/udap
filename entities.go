package main

import "encoding/json"

type Entity struct {
	Persistent
	Name        string     `json:"name" gorm:"unique"`
	Description string     `json:"description"`
	Functions   []Function `json:"functions" gorm:"many2many:entityFunction;"`
}

// root <inherits>(diagnostic)
// diagnostic <inherits>(operations)
// operations <inherits>(default)
// light.on, light.off or light.*

// url: jdfkjfdldf
// payload: on

type ModuleContext struct {
	// JSON ENV
}

type Module interface {
	Run(env ModuleContext, identifier string) (bool, error)
	Poll(env ModuleContext, identifier string) (interface{}, error)
}

type Instance struct {
	Persistent
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
	Permission  string `json:"permission"`
	environment string // Private, unless polled manually by function
}

type Object map[string]interface{}

func (i *Instance) Environment() (o Object) {
	err := json.Unmarshal([]byte(i.environment), &o)
	if err != nil {
		return nil
	}
	return o
}

type Permission struct {
	Persistent
	Name        string `json:"name"  gorm:"unique"`
	Identifier  string `json:"identifier"  gorm:"unique"`
	Description string `json:"description"`
}
