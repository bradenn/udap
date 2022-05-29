// Copyright (c) 2021 Braden Nicholson

package main

import (
	"os/exec"
	"udap/internal/core/domain"
	"udap/pkg/plugin"
)

var Module MacMeta

type MacMeta struct {
	plugin.Module
	localDisplay bool
}

func init() {
	config := plugin.Config{
		Name:        "macmeta",
		Type:        "module",
		Description: "MacOS meta interface for udap",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
}

func (v *MacMeta) createDisplaySwitch() error {

	newSwitch := &domain.Entity{
		Name:   "terminal",
		Type:   "switch",
		Module: "macmeta",
	}
	err := v.Entities.Register(newSwitch)
	if err != nil {
		return err
	}
	on := &domain.Attribute{
		Key:     "on",
		Value:   "true",
		Request: "true",
		Type:    "toggle",
		Channel: make(chan domain.Attribute),
		Order:   0,
		Entity:  newSwitch.Id,
	}

	go func() {
		for attribute := range on.Channel {
			if attribute.Request == "true" {
				err = v.displayOn()
				if err != nil {
					continue
				}
			} else {
				err = v.displayOff()
				if err != nil {
					continue
				}
			}
		}
	}()

	err = v.Attributes.Register(on)
	if err != nil {
		return err
	}
	return nil
}

func (v *MacMeta) displayOn() error {
	cmd := exec.Command("caffeinate", "-u", "-t", "1")
	err := cmd.Run()
	if err != nil {
		return err
	}
	v.localDisplay = true
	return nil
}

func (v *MacMeta) displayOff() error {
	cmd := exec.Command("pmset", "displaysleepnow")
	err := cmd.Run()
	if err != nil {
		return err
	}
	v.localDisplay = false
	return nil
}

func (v *MacMeta) Setup() (plugin.Config, error) {
	return v.Config, nil
}

func (v *MacMeta) Update() error {
	return nil
}

func (v *MacMeta) Run() error {
	err := v.displayOn()
	if err != nil {
		return err
	}
	err = v.createDisplaySwitch()
	if err != nil {
		return err
	}
	return nil
}
