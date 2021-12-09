// Copyright (c) 2021 Braden Nicholson

package main

import (
	"context"
	"fmt"
	"github.com/jaedle/golang-tplink-hs100/pkg/configuration"
	"github.com/jaedle/golang-tplink-hs100/pkg/hs100"
	"strings"
	"time"
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

	devices, err := hs100.Discover("192.168.2.0/24", configuration.Default().WithTimeout(time.Second))
	if err != nil {
		return h.Config, err
	}
	if len(devices) == 0 {
		fmt.Println("No devices")
	}
	for _, device := range devices {
		name, err := device.GetName()
		if err != nil {
			return h.Config, err
		}
		h.devices[name] = device
	}
	return h.Config, nil
}

// Update is called every cycle
func (h *HS100) Update(ctx context.Context) (err error) {
	// for name, device := range h.devices {
	// 	// setState := UpdateState(device)
	// 	// Module.UpdateState(ctx, name, &setState)
	// }
	return nil
}

// Run is called after Setup, concurrent with Update
func (h *HS100) Run() (err error) {
	for _, device := range h.devices {
		var name string
		name, err = device.GetName()
		if err != nil {
			return err
		}
		_, err = device.IsOn()
		if err != nil {
			return err
		}

		newSwitch := models.NewSwitch(strings.ToLower(name), "hs100")
		err = newSwitch.Handlers(Tx(device), Rx(device))
		if err != nil {
			return err
		}

		h.RegisterEntity(newSwitch)
	}
	return nil
}

func Tx(device *hs100.Hs100) models.Tx {
	return func(state models.State) error {
		a := models.Mono{}
		a.Unmarshal(state)
		if a.Value > 0 {
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
		a := models.Mono{}
		rm, err := device.IsOn()
		if err != nil {
			return []byte{}
		}
		if rm {
			a.Value = 1.0
		} else {
			a.Value = 0
		}
		return a.Marshal()
	}
}
