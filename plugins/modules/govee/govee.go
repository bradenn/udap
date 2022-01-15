// Copyright (c) 2021 Braden Nicholson

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gerow/go-color"
	_ "github.com/gerow/go-color"
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

func (g *Govee) getAllStates(device Device, mode string) (any, error) {
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
	for _, value := range a.Properties {
		for s, message := range value {
			switch s {
			case "powerState":
				if mode == "on" {
					var st string
					err = json.Unmarshal(message, &st)
					if err != nil {
						return nil, err
					}
					return st == "on", nil
				}
			case "colorTem":
				if mode == "cct" {
					var out int
					err = json.Unmarshal(message, &out)
					if err != nil {
						return nil, err
					}
					return out, nil
				}
			case "brightness":
				if mode == "dim" {
					var out int
					err = json.Unmarshal(message, &out)
					if err != nil {
						return nil, err
					}
					return out, nil
				}
			case "color":
				if mode == "hue" {
					var col Color
					err = json.Unmarshal(message, &col)
					if err != nil {
						return nil, err
					}
					rg := color.RGB{R: float64(col.R), G: float64(col.G), B: float64(col.B)}
					h := rg.ToHSL()
					return int(h.H), nil
				}
			}
		}
	}

	return nil, nil
}

func (g *Govee) getSingleState(device Device, mode string) (any, error) {
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

	for _, value := range a.Properties {
		for s, message := range value {
			switch s {
			case "powerState":
				if mode == "on" {
					var st string
					err = json.Unmarshal(message, &st)
					if err != nil {
						return nil, err
					}
					return st == "on", nil
				}
				continue
			case "colorTem":
				if mode == "cct" {
					var out int
					err = json.Unmarshal(message, &out)
					if err != nil {
						return nil, err
					}
					return out, nil
				}
				continue
			case "brightness":
				if mode == "dim" {
					var out int
					err = json.Unmarshal(message, &out)
					if err != nil {
						return nil, err
					}
					return out, nil
				}
				continue
			case "color":
				if mode == "hue" {
					var col Color
					err = json.Unmarshal(message, &col)
					if err != nil {
						return nil, err
					}
					rg := color.RGB{R: float64(col.R), G: float64(col.G), B: float64(col.B)}
					h := rg.ToHSL()
					return int(h.H), nil
				}
				continue
			}
		}
	}

	return nil, nil
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

func (g *Govee) setState(device Device, mode string, value any) error {
	switch mode {
	case "cct":
		cct := Cmd{
			Name:  "colorTem",
			Value: value,
		}
		err := g.control(device, cct)
		if err != nil {
			return err
		}
	case "on":
		var res bool
		err := json.Unmarshal(value.(json.RawMessage), &res)
		if err != nil {
			return err
		}
		ns := "off"
		if res {
			ns = "on"
		}
		cmd := Cmd{
			Name:  "turn",
			Value: ns,
		}
		err = g.control(device, cmd)
		if err != nil {
			return err
		}
	case "dim":
		cmd := Cmd{
			Name:  "brightness",
			Value: value,
		}
		err := g.control(device, cmd)
		if err != nil {
			return err
		}
	case "hue":
		var a int
		err := json.Unmarshal(value.(json.RawMessage), &a)
		if err != nil {
			return err
		}
		h := color.HSL{
			H: float64(a) / 360.0,
			S: 1,
			L: 0.5,
		}
		rgb := h.ToRGB()
		c := Cmd{
			Name: "color",
			Value: Color{
				R: int(rgb.R),
				B: int(rgb.G),
				G: int(rgb.B),
			},
		}
		err = g.control(device, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Govee) statePut(device Device, mode string) func(any) error {
	return func(body any) error {
		return g.setState(device, mode, body)
	}
}

func (g *Govee) stateGet(device Device, mode string) func() (any, error) {
	return func() (any, error) {
		return g.getSingleState(device, mode)
	}
}

func (g *Govee) Setup() (plugin.Config, error) {

	devices, err := g.fetchDevices()
	if err != nil {
		return g.Config, err
	}

	for _, device := range devices {

		s := models.NewSpectrum(device.DeviceName, g.Config.Name)
		for _, attribute := range s.Attributes {
			attribute.FnGet(g.stateGet(device, attribute.Key))
			attribute.FnPut(g.statePut(device, attribute.Key))
		}

		_, err = g.Entities.Register(s)
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
