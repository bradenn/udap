// Copyright (c) 2022 Braden Nicholson

package controller

import (
	"udap/internal/core/domain"
	"udap/internal/core/ports"
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
	RX            chan<- domain.Mutation
}

type CoreModule interface {
	Watch(chan domain.Mutation) error
	EmitAll() error
}

func NewController(resp chan domain.Mutation) (*Controller, error) {
	c := &Controller{
		RX: resp,
	}
	return c, nil
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
