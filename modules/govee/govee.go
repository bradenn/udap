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
	"strconv"
	"sync"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/plugin"
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
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
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

type RateLimit struct {
	Remaining string `json:"remaining"`
	Limit     string `json:"limit"`
	Reset     string `json:"reset"`
}

type Response struct {
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
	Code    int             `json:"code"`
}

var Module Govee

type Govee struct {
	plugin.Module
	devices   map[string]Device
	mutable   chan domain.Attribute
	immutable chan domain.Attribute
	done      chan bool
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

func (g *Govee) listen() {
	for {
		select {
		case attribute := <-g.mutable:
			go func() {
				err := g.setState(attribute)
				if err != nil {
					g.Err(err)
				}
			}()
		case attribute := <-g.immutable:
			g.WarnF("Attribute '%s' is immutable", attribute.Key)
		case <-g.done:
			return
		}
	}
}

func (g *Govee) Dispose() error {
	select {
	case g.done <- true:
	default:
		return nil
	}

	return nil
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
	//
	// fmt.Println(d)

	return d.Devices, nil
}

// a5:20:d4:ad:fc:08:b0:b3 H6003 Entrance
// 98:34:d4:ad:fc:0a:3f:2d H6003 Workstation
// 3d:b2:d4:ad:fc:09:38:0f H6003 Kitchen
// 6d:2f:d4:ad:fc:09:3f:25 H6003 Nightstand

func (g *Govee) sendApiDeviceRequest(id string, method string, path string, body json.RawMessage) (json.RawMessage,
	error) {

	var buf bytes.Buffer
	buf.Write(body)

	c := http.Client{}
	defer c.CloseIdleConnections()
	c.Timeout = time.Millisecond * 2000
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

	rl := RateLimit{
		Remaining: do.Header.Get("API-RateLimit-Remaining"),
		Limit:     do.Header.Get("API-RateLimit-Limit"),
		Reset:     do.Header.Get("API-RateLimit-Reset"),
	}

	marshal, err := json.Marshal(rl)
	if err != nil {
		return nil, err
	}

	err = g.Attributes.Update(id, "api", string(marshal), time.Now())
	if err != nil {
		return nil, err
	}

	_ = do.Body.Close()

	if r.Code != 200 {
		g.WarnF("API Request failed: %s", fmt.Errorf("update failed '%s'", r.Message))
		return nil, fmt.Errorf("update failed '%s'", r.Message)
	}
	return r.Data, nil
}

func (g *Govee) sendApiRequest(method string, path string, body json.RawMessage) (json.RawMessage, error) {

	var buf bytes.Buffer
	buf.Write(body)

	c := http.Client{}
	defer c.CloseIdleConnections()
	c.Timeout = time.Millisecond * 3000
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

	_ = do.Body.Close()

	if r.Code != 200 {
		g.WarnF("API Request failed: %s", fmt.Errorf("update failed '%s'", r.Message))
		return nil, fmt.Errorf("update failed '%s'", r.Message)
	}

	return r.Data, nil
}

func (g *Govee) getAllStates(device Device, id string) (interface{}, error) {
	path := fmt.Sprintf("/state?device=%s&&model=%s", device.Device, device.Model)
	stamp := time.Now()
	request, err := g.sendApiDeviceRequest(id, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	a := StateResponse{}
	err = json.Unmarshal(request, &a)
	if err != nil {
		return nil, err
	}

	for _, value := range a.Properties {
		for s, message := range value {
			switch s {
			case "online":
				var str string
				err = json.Unmarshal(message, &str)
				if err != nil {
					break
				}
				err = g.Attributes.Update(id, "online", str, stamp)
				if err != nil {
					fmt.Println(err)
					break
				}
				break
			case "powerState":
				var str string
				err = json.Unmarshal(message, &str)
				if err != nil {
					break
				}
				res := "false"
				if str == "on" {
					res = "true"
				}
				err = g.Attributes.Update(id, "on", res, stamp)
				if err != nil {
					fmt.Println(err)
					break
				}
				break
			case "colorTem":
				err = g.Attributes.Update(id, "cct", string(message), stamp)
				if err != nil {
					break
				}
				break
			case "brightness":
				parseInt, err := strconv.ParseInt(string(message), 10, 64)
				if err != nil {
					return nil, err
				}
				if parseInt == 0 {
					err = g.Attributes.Update(id, "on", "false", stamp)
					if err != nil {
						fmt.Println(err)
						break
					}
				}
				err = g.Attributes.Update(id, "dim", fmt.Sprintf("%d", parseInt), stamp)
				if err != nil {
					break
				}
				break
			case "color":
				var col Color
				err = json.Unmarshal(message, &col)
				if err != nil {
					break
				}
				rg := color.RGB{R: float64(col.R) / 255, G: float64(col.G) / 255, B: float64(col.B) / 255}
				h := rg.ToHSL()
				marshal, err := json.Marshal(int(h.H * 360))
				if err != nil {
					break
				}
				err = g.Attributes.Update(id, "hue", string(marshal), stamp)
				if err != nil {
					break
				}
				break
			default:
				break
			}
		}
	}
	return nil, nil
}

func (g *Govee) control(id string, cmd Cmd) error {
	device := g.devices[id]
	r := SetStateRequest{
		Device: device.Device,
		Model:  device.Model,
		Cmd:    cmd,
	}
	marshal, err := json.Marshal(r)
	if err != nil {
		return err
	}

	_, err = g.sendApiDeviceRequest(id, "PUT", "/control", marshal)
	if err != nil {
		return err
	}

	return nil
}

func (g *Govee) setOn(id string, b bool) error {
	ns := "off"
	if b {
		ns = "on"
	}
	cmd := Cmd{
		Name:  "turn",
		Value: ns,
	}
	err := g.control(id, cmd)
	if err != nil {
		return err
	}

	return nil

}

func (g *Govee) setCCT(id string, b int) error {
	cmd := Cmd{
		Name:  "colorTem",
		Value: b,
	}
	err := g.control(id, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (g *Govee) setLevel(id string, b int) error {
	cmd := Cmd{
		Name:  "brightness",
		Value: b,
	}
	err := g.control(id, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (g *Govee) setHue(id string, hue int) error {
	h := color.HSL{
		H: float64(hue) / 360.0,
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
	err := g.control(id, c)
	if err != nil {
		return err
	}
	return nil
}

func (g *Govee) setState(attribute domain.Attribute) error {
	// device Device, value string, mode string, id string
	id := attribute.Entity
	value := attribute.Request
	mode := attribute.Key
	switch mode {
	case "cct":
		val, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		err = g.setCCT(id, val)
		if err != nil {
			return err
		}
		break
	case "on":
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}

		err = g.setOn(id, parsed)
		if err != nil {
			return err
		}
		break
	case "dim":
		val, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		err = g.setLevel(id, val)
		if err != nil {
			return err
		}
		break
	case "hue":
		val, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		err = g.setHue(id, val)
		if err != nil {
			return err
		}
		break
	}
	err := g.Attributes.Update(id, mode, value, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (g *Govee) generateAttributes(id string) []*domain.Attribute {
	// Mutable
	on := domain.Attribute{
		Key:     "on",
		Value:   "false",
		Request: "false",
		Type:    "toggle",
		Order:   0,
		Entity:  id,
		Channel: g.mutable,
	}
	dim := domain.Attribute{
		Key:     "dim",
		Value:   "0",
		Request: "0",
		Type:    "range",
		Order:   1,
		Entity:  id,
		Channel: g.mutable,
	}
	cct := domain.Attribute{
		Key:     "cct",
		Value:   "2000",
		Request: "2000",
		Type:    "range",
		Order:   3,
		Entity:  id,
		Channel: g.mutable,
	}
	hue := domain.Attribute{
		Key:     "hue",
		Value:   "0",
		Request: "0",
		Type:    "range",
		Order:   4,
		Entity:  id,
		Channel: g.mutable,
	}
	// Immutable
	api := domain.Attribute{
		Key:     "api",
		Value:   "{remaining: 0, limit: 0, reset: 0}",
		Request: "0",
		Type:    "media",
		Order:   5,
		Entity:  id,
		Channel: g.immutable,
	}
	online := domain.Attribute{
		Key:     "online",
		Value:   "false",
		Request: "false",
		Type:    "toggle",
		Order:   6,
		Entity:  id,
		Channel: g.immutable,
	}
	return []*domain.Attribute{&on, &dim, &hue, &cct, &api, &online}
}

func (g *Govee) Setup() (plugin.Config, error) {
	g.devices = map[string]Device{}
	err := g.UpdateInterval(1000 * 60)
	if err != nil {
		return plugin.Config{}, err
	}
	return g.Config, nil
}

func (g *Govee) push() error {
	wg := sync.WaitGroup{}
	wg.Add(len(g.devices))
	for s, d := range g.devices {
		go func(device Device, id string) {
			defer wg.Done()
			_, err := g.getAllStates(device, id)
			if err != nil {
				log.Err(err)
			}
		}(d, s)
	}
	wg.Wait()

	return nil
}

func (g *Govee) Update() error {
	if g.Ready() {
		return g.push()
	}
	return nil
}

func (g *Govee) Run() error {
	err := g.InitConfig("key", "<api key>")
	if err != nil {
		return err
	}

	devices, err := g.fetchDevices()
	if err != nil {
		g.Err(err)
		return err
	}
	g.immutable = make(chan domain.Attribute)
	g.mutable = make(chan domain.Attribute)
	g.done = make(chan bool)
	go g.listen()
	for _, device := range devices {

		s := &domain.Entity{
			Name:   device.DeviceName,
			Type:   "spectrum",
			Module: g.Config.Name,
		}
		err = g.Entities.Register(s)
		if err != nil {
			return err
		}

		g.devices[s.Id] = device
		attributes := g.generateAttributes(s.Id)
		for _, attribute := range attributes {
			err = g.Attributes.Register(attribute)
			if err != nil {
				return err
			}
		}

	}
	err = g.push()
	if err != nil {
		return err
	}
	return nil
}
