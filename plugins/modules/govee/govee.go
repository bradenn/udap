// Copyright (c) 2021 Braden Nicholson

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"udap/internal/log"
	"udap/internal/models"
	"udap/pkg/plugin"
)

type Color struct {
	R int `json:"r"`
	B int `json:"b"`
	G int `json:"g"`
}

type Property struct {
	Online     bool   `json:"online,omitempty"`
	PowerState string `json:"powerState,omitempty"`
	Brightness int    `json:"brightness,omitempty"`
	Color      Color  `json:"color,omitempty"`
}

type DevicesResponse struct {
	Devices []Device `json:"devices"`
}

type StateResponse struct {
	Properties []map[string]json.RawMessage `json:"properties"`
	Device     string                       `json:"device"`
	Model      string                       `json:"model"`
}

type Cmd struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type SetStateRequest struct {
	Device string `json:"device"`
	Model  string `json:"model"`
	Cmd    `json:"cmd"`
}

type Device struct {
	Device       string   `json:"device"`
	Model        string   `json:"model"`
	DeviceName   string   `json:"deviceName"`
	Controllable bool     `json:"controllable"`
	Retrievable  bool     `json:"retrievable"`
	SupportCmds  []string `json:"supportCmds"`
	Properties   struct {
		ColorTem struct {
			Range struct {
				Min int `json:"min"`
				Max int `json:"max"`
			} `json:"range"`
		} `json:"colorTem"`
	} `json:"properties"`
}

type Response struct {
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
	Code    int             `json:"code"`
}

var Module Govee

type Govee struct {
	plugin.Module
	devices map[string]Device
}

func init() {
	config := plugin.Config{
		Name:        "govee",
		Type:        "module",
		Description: "Govee Light Controller",
		Version:     "0.1.1",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (g *Govee) fetchDevices() ([]Device, error) {
	request, err := g.sendApiRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}

	d := DevicesResponse{}
	err = json.Unmarshal(request, &d)
	if err != nil {
		return nil, err
	}

	return d.Devices, nil
}

func (g *Govee) sendApiRequest(method string, path string, body json.RawMessage) (json.RawMessage, error) {

	var buf bytes.Buffer
	buf.Write(body)
	c := http.Client{}
	p := fmt.Sprintf("https://developer-api.govee.com/v1/devices%s", path)
	request, err := http.NewRequest(method, p, &buf)
	if err != nil {
		return nil, err //
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Govee-API-Key", os.Getenv("goveeApi"))
	do, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	buf = bytes.Buffer{}
	_, err = buf.ReadFrom(do.Body)
	if err != nil {
		return nil, err
	}

	r := Response{}
	err = json.Unmarshal(buf.Bytes(), &r)
	if err != nil {
		return nil, err
	}

	if r.Code != 200 {
		return nil, fmt.Errorf("govee responded with a non 200 status code: %g", r.Message)
	}

	return r.Data, nil
}

func (g *Govee) getState(device Device) (models.State, error) {
	path := fmt.Sprintf("/state?device=%s&&model=%s", device.Device, device.Model)

	request, err := g.sendApiRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	a := StateResponse{}
	err = json.Unmarshal(request, &a)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	state := models.LightState{}
	for _, value := range a.Properties {
		for s, message := range value {

			switch s {
			case "powerState":
				var st string
				err = json.Unmarshal(message, &st)
				if err != nil {
					return nil, err
				}
				state.Power = st
			case "colorTem":
				var out int
				err = json.Unmarshal(message, &out)
				if err != nil {
					return nil, err
				}
				state.CCT = out
			case "brightness":
				var out int
				err = json.Unmarshal(message, &out)
				if err != nil {
					return nil, err
				}
				state.Level = out
			case "color":
				var col Color
				err = json.Unmarshal(message, &col)
				if err != nil {
					return nil, err
				}
				state.Red = col.R
				state.Green = col.G
				state.Blue = col.B
			}
		}
	}
	marshal, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	return marshal, nil
}

func (g *Govee) control(device Device, cmd Cmd) error {
	r := SetStateRequest{
		Device: device.Device,
		Model:  device.Model,
		Cmd:    cmd,
	}

	marshal, err := json.Marshal(r)
	if err != nil {
		return err
	}

	request, err := g.sendApiRequest("PUT", "/control", marshal)
	if err != nil {
		return err
	}

	re := Response{}

	err = json.Unmarshal(request, &re)
	if err != nil {
		return err
	}

	return nil
}

func (g *Govee) setState(device Device, state models.LightState) error {
	switch state.Mode {
	case "cct":
		cct := Cmd{
			Name:  "colorTem",
			Value: state.CCT,
		}
		err := g.control(device, cct)
		if err != nil {
			return err
		}
	case "power":
		cct := Cmd{
			Name:  "turn",
			Value: state.Power,
		}
		err := g.control(device, cct)
		if err != nil {
			return err
		}
	case "brightness":
		cct := Cmd{
			Name:  "brightness",
			Value: state.Level,
		}
		err := g.control(device, cct)
		if err != nil {
			return err
		}
	case "color":
		color := Cmd{
			Name: "color",
			Value: Color{
				R: state.Red,
				B: state.Blue,
				G: state.Green,
			},
		}
		err := g.control(device, color)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Govee) Setup() (plugin.Config, error) {

	devices, err := g.fetchDevices()
	if err != nil {
		return g.Config, err
	}

	for _, device := range devices {

		state := models.LightState{
			Power: "off",
			Red:   0,
			Green: 0,
			Blue:  0,
			Level: 0,
			CCT:   5000,
		}
		e := models.NewWildcardDevice(device.DeviceName, g.Config.Name, models.State(state.JSON()))

		err = e.Handlers(g.Tx(device), g.Rx(device))
		if err != nil {
			return plugin.Config{}, err
		}
		_, err = g.Send("entity", "register", e)
		if err != nil {
			return plugin.Config{}, err
		}

	}

	return g.Config, nil
}

func (g *Govee) Update() error {
	return nil
}

func (g *Govee) Run() error {
	return nil
}

func (g *Govee) Tx(device Device) models.Tx {
	return func(state models.State) error {
		a := models.LightState{}
		err := json.Unmarshal(state, &a)
		if err != nil {
			return err
		}
		err = g.setState(device, a)
		if err != nil {
			return err
		}
		return nil
	}
}

func (g *Govee) Rx(device Device) models.Rx {
	return func() models.State {
		state, err := g.getState(device)
		if err != nil {
			return nil
		}
		return state
	}
}

// func (s *Govee) Tx(channel int) models.Tx {
// 	return func(state models.State) error {
// 		mono := models.Mono{}
// 		err := json.Unmarshal(state, &mono)
// 		if err != nil {
// 			return err
// 		}
// 		s.stateMutex.Lock()
// 		s.state[channel] = int(mono.Value * 255.0)
// 		s.stateMutex.Unlock()
// 		return nil
// 	}
// }
//
// func (s *Govee) Rx(channel int) models.Rx {
// 	return func() models.State {
// 		mono := models.Mono{}
// 		s.stateMutex.Lock()
// 		mono.Value = float32(s.state[channel]) / 255.0
// 		s.stateMutex.Unlock()
// 		return mono.Marshal()
// 	}
// }
