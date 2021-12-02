// Copyright (c) 2021 Braden Nicholson

package main

import (
	"fmt"
	"github.com/jaedle/golang-tplink-hs100/pkg/configuration"
	"github.com/jaedle/golang-tplink-hs100/pkg/hs100"
	"udap/internal/log"
	"udap/pkg/plugin"

	"time"
)

var Plugin = HS100{
	metadata: plugin.Metadata{
		Name:        "HS100",
		Type:        "module",
		Description: "Control TP-Link HS100 Outlets",
		Version:     "1.3.1",
		Author:      "Braden Nicholson",
	}}

type HS100 struct {
	plugin.SDK
	metadata plugin.Metadata
	requests chan plugin.Request
	buffer   chan plugin.Event
}

func (h *HS100) Startup() (plugin.Metadata, error) {
	h.SDK = plugin.SDK{}
	err := Plugin.discover()
	if err != nil {
		return plugin.Metadata{}, err
	}
	return h.metadata, nil
}

func (h *HS100) Connect(events chan plugin.Event) chan plugin.Request {
	h.buffer = events
	return nil
}

func (h *HS100) Metadata() plugin.Metadata {
	return h.metadata
}

func (h *HS100) Request() chan plugin.Request {
	return nil
}

func (h *HS100) Resolve(events chan plugin.Event) {

}

func (h *HS100) Cleanup() {
	close(h.buffer)
}

func (h *HS100) Listen() {
	for request := range h.requests {
		fmt.Println(request.Operation)
	}
}

// func (h *HS100) Run(ctx module.Context, data string) (string, error) {
// 	state := SetState{}
// 	err := json.Unmarshal([]byte(data), &state)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	device := context[state.Name]
// 	if device == nil {
// 		return "", fmt.Errorf("Shit -> FAN")
// 	}
//
// 	if state.State == "off" {
// 		err = device.TurnOff()
// 		if err != nil {
// 			return "", err
// 		}
// 		err = store.Put(fmt.Sprintf("instance.%s.entity.%s.state", state.Id, state.Name), "off")
// 		if err != nil {
// 			return "", err
// 		}
// 		return "off", err
// 	} else {
// 		err = device.TurnOn()
// 		if err != nil {
// 			return "", err
// 		}
// 		err = store.Put(fmt.Sprintf("instance.%s.entity.%s.state", state.Id, state.Name), "on")
// 		if err != nil {
// 			return "", err
// 		}
// 		return "on", err
// 	}
// }

func (h *HS100) discover() (err error) {
	devices, err := hs100.Discover("192.168.2.0/24", configuration.Default().WithTimeout(time.Second))
	if err != nil {
		return err
	}
	for _, device := range devices {
		var name string
		name, err = device.GetName()
		if err != nil {
			log.Log(err.Error())
		}

		err = h.SDK.CreateOrInitEntity(name, plugin.TOGGLE, EntityHandler(device))
		if err != nil {
			continue
		}
	}
	return nil
}

func EntityHandler(device *hs100.Hs100) plugin.EntityHandler {
	return func(payload string) error {
		if payload == "on" {
			err := device.TurnOn()
			if err != nil {
				return err
			}
		} else if payload == "off" {
			err := device.TurnOff()
			if err != nil {
				return err
			}
		}
		return nil
	}
}
