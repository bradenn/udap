package main

import (
	"encoding/json"
	"fmt"
	"github.com/jaedle/golang-tplink-hs100/pkg/configuration"
	"github.com/jaedle/golang-tplink-hs100/pkg/hs100"

	"time"
	"udap/types"
)

var IModule = HS100{}

var context = map[string]*hs100.Hs100{}

type HS100 struct {
}

type Config struct {
	Subnet  string   `json:"subnet"`
	Devices []string `json:"devices"`
}

func init() {
	err := IModule.discover()
	if err != nil {
		return
	}
}

func (h *HS100) Load(agent types.Agent) error {
	err := agent.Update(struct{ Msg string }{Msg: "Hello their trucker!"})
	if err != nil {
		return err
	}
	return h.Startup(agent)
}

func (h *HS100) Startup(agent types.Agent) error {
	// err := h.discover(agent)
	//
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (h *HS100) Create(agent types.Agent) error {
	config := Config{
		Subnet:  "192.168.2.0/24",
		Devices: []string{},
	}
	err := agent.Store(config)
	if err != nil {
		return err
	}
	return nil
}

func (h *HS100) Metadata() types.Metadata {
	return types.Metadata{
		Name:        "HS100",
		Description: "Control TP-Link HS100 bulbs",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
}

func (h *HS100) discover() (err error) {
	fmt.Println("HS100:", "discovering")
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
		fmt.Println("HS100:", name)
		if err != nil {
			continue
		}
	}
	return nil
}

type SetState struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

func (h *HS100) Run(data string) (error, string) {
	state := SetState{}
	err := json.Unmarshal([]byte(data), &state)
	if err != nil {
		return err, ""
	}

	device := context[state.Name]
	if device == nil {
		return fmt.Errorf("Shit -> FAN"), ""
	}

	if state.State == "off" {
		err = device.TurnOff()
		if err != nil {
			return err, ""
		}
		return err, "off"
	} else {
		err = device.TurnOn()
		if err != nil {
			return err, ""
		}
		return err, "on"
	}
}
