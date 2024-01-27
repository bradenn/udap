// Copyright (c) 2021 Braden Nicholson

package main

import (
	"github.com/jaedle/golang-tplink-hs100/pkg/configuration"
	"github.com/jaedle/golang-tplink-hs100/pkg/hs100"
	"strconv"
	"strings"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/plugin"
)

var Module HS100

type HS100 struct {
	plugin.Module
	devices map[string]*hs100.Hs100
	names   map[string]string
	mod     int
	mux     chan domain.Attribute
}

func init() {
	config := plugin.Config{
		Name:        "hs100",
		Type:        "module",
		Description: "Control TP-Link HS100 Outlets",
		Version:     "1.7.4",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (h *HS100) findDevices() error {
	devices, err := hs100.Discover("10.0.2.0/24", configuration.Default().WithTimeout(time.Second*5))
	if err != nil {
		log.Err(err)
		return nil
	}

	for len(devices) == 0 {
		log.Event("No hs110000!!!")
		return nil
	}

	for _, device := range devices {

		name, err := device.GetName()
		if err != nil {
			return err
		}

		v, ok := h.names[strings.ToLower(name)]

		if ok {
			r, here := h.devices[v]
			if here {
				_, err = r.IsOn()
				if err == nil {
					continue
				}
			}
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
		h.names[strings.ToLower(name)] = newSwitch.Id

		on := &domain.Attribute{
			Key:     "on",
			Value:   "false",
			Request: "false",
			Order:   0,
			Type:    "toggle",
			Entity:  newSwitch.Id,
			Channel: h.mux,
		}

		err = h.Attributes.Register(on)
		if err != nil {
			return err
		}

	}
	return nil
}

func (h *HS100) muxLoop() {
	for {
		select {
		case attribute := <-h.mux:
			device, ok := h.devices[attribute.Entity]
			if !ok {
				break
			}
			err := h.put(device, attribute.Entity)(attribute.Request)
			if err != nil {
				break
			}
		}

	}
}

// Setup is called once at the launch of the module
func (h *HS100) Setup() (plugin.Config, error) {
	h.devices = map[string]*hs100.Hs100{}
	h.names = map[string]string{}
	h.mod = 0
	err := h.UpdateInterval(5000)
	if err != nil {
		return plugin.Config{}, err
	}

	return h.Config, nil
}

func (h *HS100) pull() error {
	for id, device := range h.devices {

		isOn, err := device.IsOn()
		if err != nil {
		}
		res := "false"
		if isOn {
			res = "true"
		}
		err = h.Attributes.Set(id, "on", res)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update is called every cycle
func (h *HS100) Update() error {

	h.mod = (h.mod + 1) % 48
	if h.mod == 0 {
		go func() {
			err := h.findDevices()
			if err != nil {
				h.LogF("Error: %s", err.Error())
				return
			}
		}()
	}
	return h.pull()
}

// Run is called after Setup, concurrent with Update
func (h *HS100) Run() (err error) {
	h.mux = make(chan domain.Attribute, 8)
	go h.muxLoop()
	err = h.findDevices()
	if err != nil {
		return err
	}
	return nil
}

func (h *HS100) put(device *hs100.Hs100, id string) func(s string) error {
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
			err = h.Attributes.Set(id, "on", "true")
			if err != nil {
				return err
			}

		} else {
			err = device.TurnOff()
			if err != nil {
				return err
			}
			err = h.Attributes.Set(id, "on", "false")
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
