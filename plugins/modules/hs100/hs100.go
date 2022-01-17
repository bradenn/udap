// Copyright (c) 2021 Braden Nicholson

package main

import (
	"github.com/jaedle/golang-tplink-hs100/pkg/configuration"
	"github.com/jaedle/golang-tplink-hs100/pkg/hs100"
	"strconv"
	"strings"
	"time"
	"udap/internal/log"
	"udap/internal/models"
	"udap/pkg/plugin"
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
		Version:     "1.6.2",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (h *HS100) findDevices() error {
	devices, err := hs100.Discover("10.0.1.1/24", configuration.Default().WithTimeout(time.Second*4))
	if err != nil {
		log.Err(err)
	}
	for len(devices) == 0 {
		continue
	}
	for _, device := range devices {
		name, err := device.GetName()
		if err != nil {
			return err
		}

		newSwitch := models.NewSwitch(strings.ToLower(name), "hs100")

		_, err = h.Entities.Register(newSwitch)
		if err != nil {
			return err
		}

		h.devices[newSwitch.Id] = device
		on := models.Attribute{
			Key:     "on",
			Value:   "false",
			Request: "false",
			Type:    "toggle",
			Entity:  newSwitch.Id,
		}
		on.FnGet(Rx(device))
		on.FnPut(Tx(device))
		err = h.Attributes.Register(&on)
		if err != nil {
			return err
		}

	}
	return nil
}

// Setup is called once at the launch of the module
func (h *HS100) Setup() (plugin.Config, error) {
	h.devices = map[string]*hs100.Hs100{}
	err := h.findDevices()
	if err != nil {
		return plugin.Config{}, err
	}
	return h.Config, nil
}

// Update is called every cycle
func (h *HS100) Update() error {

	for id, device := range h.devices {
		isOn, err := device.IsOn()
		if err != nil {
			return err
		}
		res := "false"
		if isOn {
			res = "true"
		}
		err = h.Attributes.Update(id, "on", res)
		if err != nil {
			return err
		}
	}

	return nil
}

// Run is called after Setup, concurrent with Update
func (h *HS100) Run() (err error) {
	for {
		err = h.Update()
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 2)
	}
}

func Tx(device *hs100.Hs100) models.FuncPut {
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

func Rx(device *hs100.Hs100) models.FuncGet {
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
