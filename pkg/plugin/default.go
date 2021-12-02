// Copyright (c) 2021 Braden Nicholson

package plugin

type Config struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // Module, Daemon, etc.
	Description string `json:"description"`
	Version     string `json:"version"`
	Author      string `json:"author"`
}

type Module struct {
	Config
	loaded bool
}

func NewModule(target *Module, config Config) {
	target.Config = config
	target.loaded = true
}
