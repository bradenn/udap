// Copyright (c) 2022 Braden Nicholson

package controller

import (
	"fmt"
	"udap/internal/bond"
	"udap/internal/pulse"
)

type Controller struct {
	Entities   *Entities
	Attributes *Attributes
	Modules    *Modules
	Endpoints  *Endpoints
	Devices    *Devices
	Zones      *Zones
	Networks   *Networks
	Users      *Users
	event      chan bond.Msg
}

func NewController() (*Controller, error) {
	c := &Controller{}
	c.Entities = LoadEntities()
	c.Modules = LoadModules()
	c.Attributes = LoadAttributes()
	c.Endpoints = LoadEndpoints()
	c.Devices = LoadDevices()
	c.Networks = LoadNetworks()
	c.Users = LoadUsers()
	c.Zones = LoadZones()
	c.Modules = LoadModules()
	return c, nil
}

func (c *Controller) Handle(msg bond.Msg) (interface{}, error) {

	pulse.LogGlobal("-> Ctrl::%s %s", msg.Target, msg.Operation)

	switch t := msg.Target; t {
	case "user":
		return c.Users.Handle(msg)
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
	case "zone":
		return c.Zones.Handle(msg)
	default:
		return nil, fmt.Errorf("unknown target '%s'", t)
	}
}

func (c *Controller) EmitAll() error {
	var err error

	err = c.Entities.EmitAll()
	if err != nil {
		return err
	}

	err = c.Attributes.EmitAll()
	if err != nil {
		return err
	}

	err = c.Networks.EmitAll()
	if err != nil {
		return err
	}

	err = c.Devices.EmitAll()
	if err != nil {
		return err
	}

	err = c.Modules.EmitAll()
	if err != nil {
		return err
	}

	err = c.Endpoints.EmitAll()
	if err != nil {
		return err
	}

	err = c.Endpoints.EmitAll()
	if err != nil {
		return err
	}

	err = c.Users.EmitAll()
	if err != nil {
		return err
	}

	err = c.Zones.EmitAll()
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) Meta(msg bond.Msg) error {
	switch t := msg.Operation; t {
	default:
		return fmt.Errorf("unknown operation '%s'", t)
	}
}
