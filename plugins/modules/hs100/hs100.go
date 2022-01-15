// Copyright (c) 2021 Braden Nicholson

package main

import (
	"fmt"
	"github.com/jaedle/golang-tplink-hs100/pkg/configuration"
	"github.com/jaedle/golang-tplink-hs100/pkg/hs100"
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

// Setup is called once at the launch of the module
func (h *HS100) Setup() (plugin.Config, error) {
	h.devices = map[string]*hs100.Hs100{}
	return h.Config, nil
}

// Update is called every cycle
func (h *HS100) Update() error {

	devices, err := hs100.Discover("10.0.1.1/24", configuration.Default().WithTimeout(time.Second*4))
	if err != nil {
		log.Err(err)
	}
	for len(devices) == 0 {
		time.Sleep(time.Second * 5)
		continue
	}
	for _, device := range devices {
		name, err := device.GetName()
		if err != nil {
			return err
		}

		if h.devices[name] != nil {
			continue
		}

		h.devices[name] = device
		newSwitch := models.NewSwitch(strings.ToLower(name), "hs100")
		for _, attribute := range newSwitch.Attributes {
			attribute.FnGet(Rx(device))
			attribute.FnPut(Tx(device))
		}

		_, err = h.Entities.Register(newSwitch)
		if err != nil {
			return err
		}
	}
	return nil
}

// Run is called after Setup, concurrent with Update
func (h *HS100) Run() (err error) {

	return h.Update()
}

func Tx(device *hs100.Hs100) func(any) error {
	return func(state any) error {
		res, ok := state.(bool)
		if !ok {
			return fmt.Errorf("invalid state type")
		}
		if res {
			err := device.TurnOn()
			if err != nil {
				return err
			}
		} else {
			err := device.TurnOff()
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func Rx(device *hs100.Hs100) func() (any, error) {
	return func() (any, error) {
		on, err := device.IsOn()
		if err != nil {
			return false, err
		}
		return on, nil
	}
}
