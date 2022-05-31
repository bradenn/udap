// Copyright (c) 2022 Braden Nicholson

package controller

import (
	"gorm.io/gorm"
	"udap/internal/bond"
	"udap/internal/core/domain"
	"udap/internal/core/modules/attribute"
	"udap/internal/core/modules/device"
	"udap/internal/core/modules/entity"
	"udap/internal/core/modules/network"
	"udap/internal/core/modules/user"
	"udap/internal/core/modules/zone"
)

type Controller struct {
	Attributes domain.AttributeService
	Devices    domain.DeviceService
	Entities   domain.EntityService
	Networks   domain.NetworkService
	Users      domain.UserService
	Zones      domain.ZoneService
	Endpoints  domain.EndpointService
	Modules    domain.ModuleService
	event      chan bond.Msg
}

func NewController(db *gorm.DB) (*Controller, error) {
	c := &Controller{}
	c.Attributes = attribute.New(db)
	c.Entities = entity.New(db)
	c.Devices = device.New(db)
	c.Networks = network.New(db)
	c.Users = user.New(db)
	c.Zones = zone.New(db)

	return c, nil
}

func (c *Controller) Listen(resp chan domain.Mutation) {

	err := c.Modules.Watch(resp)
	if err != nil {
		return
	}

	err = c.Endpoints.Watch(resp)
	if err != nil {
		return
	}

	err = c.Entities.Watch(resp)
	if err != nil {
		return
	}

	err = c.Attributes.Watch(resp)
	if err != nil {
		return
	}

}

func (c *Controller) EmitAll() error {

	err := c.Entities.EmitAll()
	if err != nil {
		return err
	}

	err = c.Attributes.EmitAll()
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

	return nil
}
