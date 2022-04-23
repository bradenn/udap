// Copyright (c) 2021 Braden Nicholson

package plugin

import (
	"fmt"
	"plugin"
	"udap/internal/bond"
	"udap/internal/controller"
)

// UdapPlugin defines the functions of a plugin's exported variable
type UdapPlugin interface {
	// Setup is called when the plugin is first loaded
	Setup() (Config, error)
	// Connect c
	Connect(*controller.Controller, *bond.Bond) error
	// Run provides module-relevant data
	Run() error
	// Update is used to assign channels
	Update() error
}

// Load attempts to load the plugin from a given path
func Load(path string) (pl UdapPlugin, err error) {
	// Attempt to open that plugin
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("plugin at '%s' failed to mount: %s", path, err.Error())
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
