// Copyright (c) 2022 Braden Nicholson

package main

import (
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
)

var Module Acap

type Acap struct {
	plugin.Module
	entityId string
	receiver chan domain.Attribute
}

func init() {
	config := plugin.Config{
		Name:        "acap",
		Type:        "module",
		Description: "Automated Casualty Avoidance Protocol",
		Version:     "0.1.0 beta",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (a *Acap) Setup() (plugin.Config, error) {
	err := a.UpdateInterval(time.Minute * 5)
	if err != nil {
		return plugin.Config{}, err
	}
	return a.Config, nil
}

func (a *Acap) Run() error {

	a.receiver = make(chan domain.Attribute, 1)

	entity := &domain.Entity{
		Name:   "acap",
		Type:   "media",
		Module: a.Config.Name,
	}

	err := a.Entities.Register(entity)
	if err != nil {
		return err
	}

	a.entityId = entity.Id

	return nil
}

func (a *Acap) Update() error {
	if a.Ready() {
	}
	return nil
}
