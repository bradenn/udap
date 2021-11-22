// Copyright (c) 2021 Braden Nicholson

package main

import (
	"encoding/json"
	"fmt"
	"github.com/jaedle/golang-tplink-hs100/pkg/configuration"
	"github.com/jaedle/golang-tplink-hs100/pkg/hs100"
	"time"
	"udap/module"
	"udap/udap/store"
)

var Plugin = HS100{Metadata: module.Metadata{
	Name:        "HS100",
	Description: "Control TP-Link HS100 bulbs",
	Version:     "0.0.1",
	Author:      "Braden Nicholson",
}}

var context = map[string]*hs100.Hs100{}

type HS100 struct {
	Metadata module.Metadata
}

func (h *HS100) Startup() (module.Metadata, error) {
	err := Plugin.discover()
	if err != nil {
		return module.Metadata{}, err
	}
	return h.Metadata, nil
}

func (h *HS100) Default() interface{} {
	return Config{}
}

func (h *HS100) Load(ctx module.Context) error {
	ctx.Send(fmt.Sprintf("Hello, from '%s.'", ctx.Id()))
	return nil
}

func (h *HS100) Update(ctx module.Context) error {
	panic("implement me")
}

func (h *HS100) Run(ctx module.Context, data string) (string, error) {
	state := SetState{}
	err := json.Unmarshal([]byte(data), &state)
	if err != nil {
		return "", err
	}

	device := context[state.Name]
	if device == nil {
		return "", fmt.Errorf("Shit -> FAN")
	}

	if state.State == "off" {
		err = device.TurnOff()
		if err != nil {
			return "", err
		}
		err = store.Put(fmt.Sprintf("instance.%s.entity.%s.state", state.Id, state.Name), "off")
		if err != nil {
			return "", err
		}
		return "off", err
	} else {
		err = device.TurnOn()
		if err != nil {
			return "", err
		}
		err = store.Put(fmt.Sprintf("instance.%s.entity.%s.state", state.Id, state.Name), "on")
		if err != nil {
			return "", err
		}
		return "on", err
	}
}

type Config struct {
	Subnet  string   `json:"subnet"`
	Devices []string `json:"devices"`
}

func (h *HS100) Configure(ctx module.Context) error {
	// config := Config{
	// 	Subnet:  "192.168.2.0/24",
	// 	Devices: []string{},
	// }
	//
	//
	return nil
}

func (h *HS100) discover() (err error) {
	devices, err := hs100.Discover("192.168.2.0/24", configuration.Default().WithTimeout(time.Second))
	if err != nil {
		return err
	}
	for _, device := range devices {
		name, err := device.GetName()
		if err != nil {
			continue
		}
		context[name] = device
		if err != nil {
			continue
		}
	}
	return nil
}

type SetState struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}
