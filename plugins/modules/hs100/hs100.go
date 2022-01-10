// Copyright (c) 2021 Braden Nicholson

package main

import (
	"encoding/json"
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
		Version:     "1.3.1",
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
func (h *HS100) Update() (err error) {
	for {
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
			err = newSwitch.Handlers(Tx(device), Rx(device))
			if err != nil {
				continue
			}
			_, err = h.Send("entity", "register", newSwitch)
			if err != nil {
				return err
			}
		}
		time.Sleep(time.Second * 30)
	}
}

// Run is called after Setup, concurrent with Update
func (h *HS100) Run() (err error) {

	return h.Update()
}

func Tx(device *hs100.Hs100) models.Tx {
	return func(state models.State) error {
		a := models.LightState{}
		err := json.Unmarshal(state, &a)
		if err != nil {
			return err
		}

		if a.Power == "on" {
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

func Rx(device *hs100.Hs100) models.Rx {
	return func() models.State {
		a := models.LightState{}
		rm, err := device.IsOn()
		if err != nil {
			return []byte{}
		}
		if rm {
			a.Power = "on"
		} else {
			a.Power = "off"
		}
		marshal, err := json.Marshal(a)
		if err != nil {
			return nil
		}
		return marshal
	}
}
