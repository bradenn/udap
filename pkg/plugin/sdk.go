// Copyright (c) 2021 Braden Nicholson

package plugin

import (
	"fmt"
	"plugin"
)

type EntityType int

const (
	TOGGLE = iota
	DIMMER
)

type PConfig struct {
	Name    string `json:"name"`
	Usage   string `json:"usage"`
	Type    string `json:"type"`
	Author  string `json:"author"`
	Version string `json:"version"`
}

type Entity struct {
	Usage string `json:"usage"`
}

type SDK struct {
	config PConfig
}

func (s *SDK) Config(config PConfig) {
	s.config = config
}

// Load attempts to load the plugin from a given path
func Load(path string) (pl UdapPlugin, err error) {
	// Attempt to open that plugin
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open plugin at path '%s': %s", path, err.Error())
	}
	// Attempt to access the Plugin variable to interface with the code
	lookup, err := p.Lookup("Module")
	if err != nil {
		return nil, fmt.Errorf("plugin '%s' does not define a Plugin interface", path)
	}
	pl = lookup.(UdapPlugin)
	// Return no errors
	return pl, nil
}
