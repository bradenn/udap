// Copyright (c) 2021 Braden Nicholson

package main

import (
	"github.com/jaedle/golang-tplink-hs100/pkg/configuration"
	"github.com/jaedle/golang-tplink-hs100/pkg/hs100"
	"strconv"
	"strings"
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
)

var Module HS100

type HS100 struct {
	plugin.Module
	devices map[string]*hs100.Hs100
}

func init() {
	config := plugin.Config{
		Name:        "hs100",
		Type:        "module",
		Description: "Control TP-Link HS100 Outlets",
		Version:     "1.7.3",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (h *HS100) findDevices() error {
	devices, err := hs100.Discover("10.0.1.1/24", configuration.Default().WithTimeout(time.Second*5))
	if err != nil {
		return nil
	}

	for len(devices) == 0 {
		return nil
	}

	for _, device := range devices {
		name, err := device.GetName()
		if err != nil {
			return err
		}

		newSwitch := domain.Entity{
			Name:   strings.ToLower(name),
			Type:   "switch",
			Module: "hs100",
		}
		err = h.Entities.Register(&newSwitch)
		if err != nil {
			return err
		}

		h.devices[newSwitch.Id] = device
		channel := make(chan domain.Attribute)
		on := &domain.Attribute{
			Key:     "on",
			Value:   "false",
			Request: "false",
			Order:   0,
			Type:    "toggle",
			Entity:  newSwitch.Id,
			Channel: channel,
		}

		go func() {
			for attribute := range channel {
				err = h.put(device)(attribute.Request)
				if err != nil {
					return
				}
			}
		}()

		err = h.Attributes.Register(on)
		if err != nil {
			return err
		}

	}
	return nil
}

// Setup is called once at the launch of the module
func (h *HS100) Setup() (plugin.Config, error) {
	h.devices = map[string]*hs100.Hs100{}

	return h.Config, nil
}

func (h *HS100) pull() error {
	for id, device := range h.devices {

		isOn, err := device.IsOn()
		if err != nil {
			return nil
		}
		res := "false"
		if isOn {
			res = "true"
		}
		err = h.Attributes.Update(id, "on", res, time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}

// Update is called every cycle
func (h *HS100) Update() error {
	if time.Since(h.Module.LastUpdate) >= time.Second*2 {
		h.Module.LastUpdate = time.Now()
		return h.pull()
	}
	return nil
}

// Run is called after Setup, concurrent with Update
func (h *HS100) Run() (err error) {
	err = h.findDevices()
	if err != nil {
		return err
	}
	return nil
}

func (h *HS100) put(device *hs100.Hs100) func(s string) error {
	return func(s string) error {

		parseBool, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		if parseBool {
			err = device.TurnOn()
			if err != nil {
				return err
			}
		} else {
			err = device.TurnOff()
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (h *HS100) get(device *hs100.Hs100) func() (string, error) {
	return func() (string, error) {
		on, err := device.IsOn()
		if err != nil {
			return "false", err
		}
		if on {
			return "true", nil
		}
		return "false", nil
	}
}
