// Copyright (c) 2022 Braden Nicholson

package controller

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/services/attribute"
	"udap/internal/core/services/device"
	"udap/internal/core/services/entity"
	"udap/internal/core/services/logs"
	"udap/internal/core/services/network"
	"udap/internal/core/services/notification"
	"udap/internal/core/services/user"
	"udap/internal/core/services/zone"
)

type Controller struct {
	Attributes    domain.AttributeService
	Devices       domain.DeviceService
	Entities      domain.EntityService
	Networks      domain.NetworkService
	Logs          domain.LogService
	Notifications domain.NotificationService
	Users         domain.UserService
	Zones         domain.ZoneService
	Endpoints     domain.EndpointService
	Modules       domain.ModuleService
	Macros        domain.MacroService
}

type CoreModule interface {
	Watch(chan domain.Mutation) error
	EmitAll() error
}

func NewController(db *gorm.DB) (*Controller, error) {
	c := &Controller{}
	c.Attributes = attribute.New(db)
	c.Entities = entity.New(db)
	c.Devices = device.New(db)
	c.Networks = network.New(db)
	c.Users = user.New(db)
	c.Notifications = notification.New(db)
	c.Zones = zone.New(db)
	c.Logs = logs.New()

	return c, nil
}

func (c *Controller) WatchAll(resp chan domain.Mutation) {

	err := c.Attributes.Watch(resp)
	if err != nil {
		return
	}

	err = c.Entities.Watch(resp)
	if err != nil {
		return
	}

	err = c.Modules.Watch(resp)
	if err != nil {
		return
	}

	err = c.Endpoints.Watch(resp)
	if err != nil {
		return
	}

	err = c.Devices.Watch(resp)
	if err != nil {
		return
	}

	err = c.Networks.Watch(resp)
	if err != nil {
		return
	}

	err = c.Zones.Watch(resp)
	if err != nil {
		return
	}

	err = c.Notifications.Watch(resp)
	if err != nil {
		return
	}

	err = c.Users.Watch(resp)
	if err != nil {
		return
	}

	err = c.Logs.Watch(resp)
	if err != nil {
		return
	}

	err = c.Macros.Watch(resp)
	if err != nil {
		return
	}

}

func (c *Controller) EmitAll() error {

	err := c.Entities.EmitAll()
	if err != nil {
		return err
	}

	err = c.Zones.EmitAll()
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

	err = c.Users.EmitAll()
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

	err = c.Notifications.EmitAll()
	if err != nil {
		return err
	}

	err = c.Logs.EmitAll()
	if err != nil {
		return err
	}

	err = c.Macros.EmitAll()
	if err != nil {
		return err
	}

	return nil
}
