// Copyright (c) 2021 Braden Nicholson

package main

import (
	"context"
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
		Version:     "1.3.1",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func cmp(bl bool, a string, b string) string {
	if bl {
		return a
	}
	return b
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
		name, err := device.GetName()
		if err != nil {
			return err
		}
		on, err := device.IsOn()
		if err != nil {
			return err
		}

		newSwitch := models.NewSwitch(strings.ToLower(name), "hs100")

		err = newSwitch.Emplace()
		if err != nil {
			return err
		}

		log.Log("Registering")
		sh := HandleState(device)
		gs := GetState(device)
		newSwitch.SetStateHandler(sh)
		newSwitch.SetStateReceiver(gs)
		h.RegisterEntity(newSwitch)

		if on {
		} else {
			newSwitch.State = false
		}

	}
	return nil
}

func HandleState(device *hs100.Hs100) func(state interface{}) error {
	return func(state interface{}) error {
		if state.(bool) {
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

func GetState(device *hs100.Hs100) func() interface{} {
	return func() interface{} {
		rm, err := device.IsOn()
		if err != nil {
			return false
		}
		return rm
	}
}
