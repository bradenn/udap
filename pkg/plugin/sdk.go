// Copyright (c) 2021 Braden Nicholson

package plugin

import (
	"fmt"
	"math"
	"plugin"
	"strings"
	"udap/internal/cache"
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

type EnergyState struct {
	State float64
}

// ToToggle will return a boolean interpretation of the state
func (e EnergyState) ToToggle() bool {
	return math.Round(e.State) >= 0.5
}

// ToDimmer will return am 8-bit value between 0 and 256
func (e EnergyState) ToDimmer() int {
	return Map(e.State, 0, 1, 0, 256)
}

func Map(value float64, i1 float64, i2 float64, o1 float64, o2 float64) int {
	return int(math.Round(o1 + (o2-o1)*((value-i1)/(i2-i1))))
}

type EntityHandler *func(payload string) error

func (s *SDK) CreateOrInitEntity(name string, entityType EntityType, handleEntity EntityHandler) error {
	cache.WatchFn(fmt.Sprintf("entities.%s.state", strings.ToLower(name)), handleEntity)
	return nil
}

func (s *SDK) RegisterEntity(config PConfig) {
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