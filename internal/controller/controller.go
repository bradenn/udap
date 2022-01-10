// Copyright (c) 2022 Braden Nicholson

package controller

import (
	"fmt"
	"udap/internal/bond"
	"udap/internal/models"
)

type Controller struct {
	Entities  *Entities
	Modules   *Modules
	Endpoints *Endpoints
	Devices   *Devices
	Networks  *Networks
	event     chan bond.Msg
}

func NewController() (*Controller, error) {
	c := &Controller{}
	c.Entities = LoadEntities()
	c.Modules = LoadModules()
	c.Endpoints = LoadEndpoints()
	c.Devices = LoadDevices()
	c.Networks = LoadNetworks()
	return c, nil
}

func (c *Controller) Handle(msg bond.Msg) (any, error) {
	switch t := msg.Target; t {
	case "entity":
		return c.Entities.Handle(msg)
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

type Metadata struct {
	Endpoint  models.Endpoint   `json:"endpoint"`
	Endpoints []models.Endpoint `json:"endpoints"`
	Devices   []models.Device   `json:"devices"`
	Entities  []models.Entity   `json:"entities"`
	Networks  []models.Network  `json:"networks"`
	Logs      []models.Log      `json:"logs"`
}
