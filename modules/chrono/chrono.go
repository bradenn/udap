// Copyright (c) 2023 Braden Nicholson

package chrono

import (
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
)

var Module Chrono

type Chrono struct {
	plugin.Module
	request chan domain.Attribute
}

type entry struct {
	Name string    `json:"name"`
	Time time.Time `json:"time"`
	Id   string    `json:"id"`
}

func init() {

	config := plugin.Config{
		Name:        "chrono",
		Type:        "module",
		Description: "A date and time module",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (c *Chrono) mux() {
	for attribute := range c.request {
		c.handleRequest(attribute)
	}
}

func (c *Chrono) handleRequest(attribute domain.Attribute) {

}

func (c *Chrono) Setup() (plugin.Config, error) {
	return Module.Config, nil
}

func (c *Chrono) Update() error {
	return nil
}

func (c *Chrono) Run() error {
	c.request = make(chan domain.Attribute)
	return nil
}

func (c *Chrono) Dispose() error {
	close(c.request)
	return nil
}
