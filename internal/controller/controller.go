// Copyright (c) 2022 Braden Nicholson

package controller

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
	"udap/internal/core/services"
)

type Controller struct {
	Attributes    ports.AttributeService
	Devices       ports.DeviceService
	Entities      ports.EntityService
	Networks      ports.NetworkService
	Logs          ports.LogService
	Notifications ports.NotificationService
	Users         ports.UserService
	Zones         ports.ZoneService
	Endpoints     ports.EndpointService
	Modules       ports.ModuleService
	Macros        ports.MacroService
	Triggers      ports.TriggerService
	SubRoutines   ports.SubRoutineService
}

type CoreModule interface {
	Watch(chan domain.Mutation) error
	EmitAll() error
}

func NewController(db *gorm.DB) (*Controller, error) {
	c := &Controller{}
	c.Entities = services.NewEntityService(db)
	c.Devices = services.NewDeviceService(db)
	c.Networks = services.NewNetworkService(db)
	c.Users = services.NewUserService(db)
	c.Notifications = services.NewNotificationService(db)
	c.Zones = services.NewZoneService(db)
	c.Logs = services.NewLogService()

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

	err = c.Triggers.Watch(resp)
	if err != nil {
		return
	}

	err = c.SubRoutines.Watch(resp)
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

	err = c.Triggers.EmitAll()
	if err != nil {
		return err
	}

	err = c.SubRoutines.EmitAll()
	if err != nil {
		return err
	}

	return nil
}
