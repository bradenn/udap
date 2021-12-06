// Copyright (c) 2021 Braden Nicholson

package plugin

import (
	"context"
)

// Metadata describes a plugin
type Metadata struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // Module, Daemon, etc.
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
}

// Event defines the parameters of an event channel request
type Event struct {
	Type      string      `json:"type"`      // Module, Daemon, Agent, etc.
	Operation string      `json:"operation"` // state
	Body      interface{} `json:"body"`      // {state: "on"}
}

// Request contains the structure of a request channel payload
type Request struct {
	Id        string `json:"id"`        // InstanceId
	Operation string `json:"operation"` // Update, Run, Poll
	Body      string `json:"body"`
}

// Plugin defines the functions of a plugin's exported variable
type Plugin interface {
	// Startup is called when the plugin is first loaded
	Startup() (Metadata, error)
	// Connect is used to assign channels
	Connect(chan Event) chan Request
}

// Device d
type Device interface {
	Name() string
}

type UdapInterface interface {
	RegisterEntity(entity *Entity)
}

// UdapPlugin defines the functions of a plugin's exported variable
type UdapPlugin interface {
	// Setup is called when the plugin is first loaded
	Setup() (Config, error)
	// Connect c
	Connect(*chan Event) (chan Request, error)
	// Run provides module-relevant data
	Run() error
	// Update is used to assign channels
	Update(ctx context.Context) error
}
