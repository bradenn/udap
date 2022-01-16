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
	"sync"
	"time"
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
	devices map[string]*Device
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

func (g *Govee) getAllStates(device Device, id string) (any, error) {
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
				var st string
				err = json.Unmarshal(message, &st)
				if err != nil {
					continue
				}
				val := "false"
				if st == "on" {
					val = "true"
				}
				err = g.Attributes.Put(id, "on", val)
				if err != nil {
					log.Err(err)
				}

				continue
			case "colorTem":
				err = g.Attributes.Put(id, "cct", string(message))
				if err != nil {
					log.Err(err)
				}
				continue
			case "brightness":
				err = g.Attributes.Put(id, "dim", string(message))
				if err != nil {
					log.Err(err)
				}
				continue
			case "color":
				var col Color
				err = json.Unmarshal(message, &col)
				if err != nil {
					continue
				}
				rg := color.RGB{R: float64(col.R) / 255, G: float64(col.G) / 255, B: float64(col.B) / 255}
				h := rg.ToHSL()
				marshal, err := json.Marshal(int(h.H * 360))
				if err != nil {
					continue
				}
				err = g.Attributes.Put(id, "hue", string(marshal))
				if err != nil {
					log.Err(err)
				}
			}
		}
	}

	return nil, nil
}

func (g *Govee) getSingleState(device Device, att models.Attribute, mode string) (string, error) {
	path := fmt.Sprintf("/state?device=%s&&model=%s", device.Device, device.Model)

	request, err := g.sendApiRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	a := StateResponse{}
	err = json.Unmarshal(request, &a)
	if err != nil {
		log.Err(err)
		return "", err
	}

	for _, value := range a.Properties {
		for s, message := range value {
			switch s {
			case "powerState":
				if mode == "on" {
					var st string
					err = json.Unmarshal(message, &st)
					if err != nil {
						return "", err
					}
					if st == "on" {
						return "true", nil
					} else {
						return "false", nil
					}
				}
				break
			case "colorTem":
				if mode == "cct" {
					var out int
					err = json.Unmarshal(message, &out)
					if err != nil {
						return "", err
					}
					return fmt.Sprintf("%d", out), nil
				}
				break
			case "brightness":
				if mode == "dim" {
					var out int
					err = json.Unmarshal(message, &out)
					if err != nil {
						return "", err
					}
					return fmt.Sprintf("%d", out), nil
				}
				break
			case "color":
				if mode == "hue" {
					var col Color
					err = json.Unmarshal(message, &col)
					if err != nil {
						return "", err
					}
					rg := color.RGB{R: float64(col.R) / 255, G: float64(col.G) / 255, B: float64(col.B) / 255}
					h := rg.ToHSL()
					return fmt.Sprintf("%d", int(h.H)*360), nil
				}
				break
			}
		}
	}

	return "", nil
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

func (g *Govee) setState(device Device, a models.Attribute) error {
	switch a.Key {
	case "cct":
		cct := Cmd{
			Name:  "colorTem",
			Value: a.AsInt(),
		}
		err := g.control(device, cct)
		if err != nil {
			return err
		}
		break
	case "on":
		ns := "off"
		if a.AsBool() {
			ns = "on"
		}
		cmd := Cmd{
			Name:  "turn",
			Value: ns,
		}
		err := g.control(device, cmd)
		if err != nil {
			return err
		}
		break
	case "dim":
		cmd := Cmd{
			Name:  "brightness",
			Value: a.AsInt(),
		}
		err := g.control(device, cmd)
		if err != nil {
			return err
		}
		break
	case "hue":

		h := color.HSL{
			H: float64(a.AsInt()) / 360.0,
			S: 1,
			L: 0.5,
		}
		rgb := h.ToRGB()
		c := Cmd{
			Name: "color",
			Value: Color{
				R: int(rgb.R * 255),
				G: int(rgb.G * 255),
				B: int(rgb.B * 255),
			},
		}
		err := g.control(device, c)
		if err != nil {
			return err
		}
		break
	}

	return nil
}

func (g *Govee) statePut(device Device, mode string) models.FuncPut {
	return func(a models.Attribute) error {
		return g.setState(device, a)
	}
}

func (g *Govee) stateGet(device Device, mode string) models.FuncGet {
	return func(a models.Attribute) (string, error) {
		return g.getSingleState(device, a, mode)
	}
}

func GenerateAttributes(id string) []*models.Attribute {
	on := models.Attribute{
		Key:    "on",
		Value:  "false",
		Type:   "toggle",
		Entity: id,
	}
	dim := models.Attribute{
		Key:    "dim",
		Value:  "0",
		Type:   "range",
		Entity: id,
	}
	hue := models.Attribute{
		Key:    "hue",
		Value:  "0",
		Type:   "range",
		Entity: id,
	}
	cct := models.Attribute{
		Key:    "cct",
		Value:  "2000",
		Type:   "range",
		Entity: id,
	}
	return []*models.Attribute{&on, &dim, &hue, &cct}
}

func (g *Govee) Setup() (plugin.Config, error) {
	g.devices = map[string]*Device{}
	devices, err := g.fetchDevices()
	if err != nil {
		return g.Config, err
	}

	for _, device := range devices {

		s := models.NewSpectrum(device.DeviceName, g.Config.Name)
		_, err = g.Entities.Register(s)
		if err != nil {
			return plugin.Config{}, err
		}

		attributes := GenerateAttributes(s.Id)
		for _, attribute := range attributes {
			attribute.FnGet(g.stateGet(device, attribute.Key))
			attribute.FnPut(g.statePut(device, attribute.Key))
			err = g.Attributes.Register(attribute)
			if err != nil {
				return plugin.Config{}, err
			}
		}

		if err != nil {
			return plugin.Config{}, err
		}
		g.devices[s.Id] = &device
	}

	return g.Config, nil
}

func (g *Govee) Update() error {
	wg := sync.WaitGroup{}
	wg.Add(len(g.devices))
	for s, d := range g.devices {
		go func(device *Device, id string) {
			defer wg.Done()
			_, err := g.getAllStates(*device, id)
			if err != nil {
				log.Err(err)
			}
		}(d, s)
	}
	wg.Wait()
	return nil
}

func (g *Govee) Run() error {
	for {
		err := g.Update()
		if err != nil {
			log.Err(err)
		}
		time.Sleep(5 * time.Second)
	}
}
