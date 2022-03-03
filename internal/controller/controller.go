// Copyright (c) 2022 Braden Nicholson

package controller

import (
	"fmt"
	"udap/internal/bond"
)

type Controller struct {
	Entities       *Entities
	Attributes     *Attributes
	Modules        *Modules
	Endpoints      *Endpoints
	Devices        *Devices
	Networks       *Networks
	UserController *UserController
	event          chan bond.Msg
}

func NewController() (*Controller, error) {
	c := &Controller{}
	c.Entities = LoadEntities()
	c.Modules = LoadModules()
	c.Attributes = LoadAttributes()
	c.Endpoints = LoadEndpoints()
	c.Devices = LoadDevices()
	c.Networks = LoadNetworks()
	return c, nil
}

func (c *Controller) Handle(msg bond.Msg) (interface{}, error) {
	switch t := msg.Target; t {
	case "user":
		return c.UserController.Handle(msg)
	case "entity":
		return c.Entities.Handle(msg)
	case "attribute":
		return c.Attributes.Handle(msg)
	case "module":
		return c.Modules.Handle(msg)
	case "endpoint":
		return c.Endpoints.Handle(msg)
	case "device":
		return c.Devices.Handle(msg)
	case "network":
		return c.Networks.Handle(msg)
	default:
		return nil, fmt.Errorf("unknown target '%s'", t)
	}
}

func (c *Controller) Meta(msg bond.Msg) error {
	switch t := msg.Operation; t {
	default:
		return fmt.Errorf("unknown operation '%s'", t)
	}
}
