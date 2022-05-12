// Copyright (c) 2021 Braden Nicholson

package main

import (
	"os/exec"
	"udap/internal/models"
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
	newSwitch := models.NewSwitch("terminal", "macmeta")
	_, err := v.Entities.Register(newSwitch)
	if err != nil {
		return err
	}
	on := &models.Attribute{
		Key:     "on",
		Value:   "true",
		Request: "true",
		Type:    "toggle",
		Order:   0,
		Entity:  newSwitch.Id,
	}

	on.FnGet(func() (string, error) {
		if v.localDisplay {
			return "true", nil
		} else {
			return "false", nil
		}
	})

	on.FnPut(func(value string) error {
		if value == "true" {
			err = v.displayOn()
			if err != nil {
				return err
			}
		} else {
			err = v.displayOff()
			if err != nil {
				return err
			}
		}
		return nil
	})

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
