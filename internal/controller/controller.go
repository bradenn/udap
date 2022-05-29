// Copyright (c) 2022 Braden Nicholson

package controller

import (
	"fmt"
	"gorm.io/gorm"
	"udap/internal/bond"
	"udap/internal/core/domain"
	"udap/internal/modules/attribute"
	"udap/internal/modules/device"
	"udap/internal/modules/endpoint"
	"udap/internal/modules/entity"
	"udap/internal/modules/network"
	"udap/internal/modules/user"
	"udap/internal/modules/zone"
	"udap/internal/pulse"
)

type Controller struct {
	Attributes domain.AttributeService
	Devices    domain.DeviceService
	Endpoints  domain.EndpointService
	Entities   domain.EntityService
	Networks   domain.NetworkService
	Users      domain.UserService
	Zones      domain.ZoneService
	event      chan bond.Msg
}

func NewController(db *gorm.DB) (*Controller, error) {
	c := &Controller{}

	c.Attributes = attribute.New(db)
	c.Devices = device.New(db)
	c.Endpoints = endpoint.New(db)
	c.Entities = entity.New(db)
	c.Networks = network.New(db)
	c.Users = user.New(db)
	c.Zones = zone.New(db)

	return c, nil
}

func (c *Controller) Handle(msg bond.Msg) (interface{}, error) {

	pulse.LogGlobal("-> Ctrl::%s %s", msg.Target, msg.Operation)

	// switch t := msg.Target; t {
	// case "attribute":
	// 	return c.Attributes.Handle(msg)
	// case "device":
	// 	return c.Devices.Handle(msg)
	// case "network":
	// 	return c.Networks.Handle(msg)
	// case "zone":
	// 	return c.Zones.Handle(msg)
	// default:
	// 	return nil, fmt.Errorf("unknown target '%s'", t)
	// }
	return nil, nil
}

func (c *Controller) EmitAll() error {
	// var err error
	// err = c.Attributes.EmitAll()
	// if err != nil {
	// 	return err
	// }
	//
	// err = c.Networks.EmitAll()
	// if err != nil {
	// 	return err
	// }
	//
	// err = c.Devices.EmitAll()
	// if err != nil {
	// 	return err
	// }
	//
	// err = c.Zones.EmitAll()
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (c *Controller) Meta(msg bond.Msg) error {
	switch t := msg.Operation; t {
	default:
		return fmt.Errorf("unknown operation '%s'", t)
	}
}
