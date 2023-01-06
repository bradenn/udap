// Copyright (c) 2021 Braden Nicholson

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"os"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/plugin"
)

var Module Sentry

type Sentry struct {
	plugin.Module
	eId             string
	position        Position
	beam            Beam
	positionChannel chan domain.Attribute
	beamChannel     chan domain.Attribute
	done            chan bool
	session         *websocket.Conn
}

type Beam struct {
	Target string  `json:"target"`
	Active int     `json:"active"`
	Power  float64 `json:"power"`
}

func (b *Beam) Marshal() string {
	marshal, err := json.Marshal(b)
	if err != nil {
		return ""
	}
	return string(marshal)
}

type Position struct {
	Pan  int `json:"pan"`
	Tilt int `json:"tilt"`
}

func (p *Position) Marshal() string {
	marshal, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(marshal)
}

const sentryUrl = "10.0.1.76"

func init() {
	config := plugin.Config{
		Name:        "sentry",
		Type:        "module",
		Description: "Sentry with LASERS",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
	Module.eId = ""
	Module.session = nil
	Module.Config = config
}

func (v *Sentry) connect() {
	u := url.URL{Scheme: "ws", Host: sentryUrl, Path: "/ws"}

	var err error
	v.session, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		v.Err(err)
		return
	}

	go func() {
		for {
			if v.session == nil {
				return
			}
			_, _, err = v.session.ReadMessage()
			if err != nil {
				v.session = nil
				return
			}
		}
	}()
}

func (v *Sentry) Setup() (plugin.Config, error) {
	err := v.UpdateInterval(2000)
	if err != nil {
		return plugin.Config{}, err
	}
	return v.Config, nil
}

type Status struct {
	System struct {
		Mac string `json:"mac"`
	} `json:"system"`
	Servos struct {
		Pan  int `json:"pan"`
		Tilt int `json:"tilt"`
	} `json:"servos"`
	Beams struct {
		Primary struct {
			Active bool    `json:"active"`
			Power  float64 `json:"power"`
		} `json:"primary"`
		Secondary struct {
			Active bool    `json:"active"`
			Power  float64 `json:"power"`
		} `json:"secondary"`
	} `json:"beams"`
}

type SetPosition struct {
	Pan   int    `json:"pan"`
	Tilt  int    `json:"tilt"`
	Token string `json:"token"`
}

func mapRange(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

func (v *Sentry) requestPosition(position SetPosition) error {
	//marshal, err := json.Marshal(position)
	//if err != nil {
	//	return err
	//}
	//reader := bytes.NewReader(marshal)
	//
	//
	//
	//client := http.Client{}
	//client.Timeout = time.Millisecond * 250
	//defer client.CloseIdleConnections()
	//resp, err := client.Post(fmt.Sprintf("http://%s/position", sentryUrl), "application/json", reader)
	//if err != nil {
	//	return err
	//}
	//
	//
	//var buf bytes.Buffer
	//_, err = buf.ReadFrom(resp.Body)
	//if err != nil {
	//	return err
	//}
	//_ = resp.Body.Close()

	if v.session == nil {
		v.connect()
	}

	err := v.session.WriteJSON(position)
	if err != nil {
		return err
	}

	return nil
}

func (v *Sentry) requestBeam(beam Beam) error {
	beam.Target = "primary"
	marshal, err := json.Marshal(beam)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(marshal)

	v.beam = beam

	client := http.Client{}
	client.Timeout = time.Millisecond * 250
	defer client.CloseIdleConnections()
	resp, err := client.Post(fmt.Sprintf("http://%s/beam", sentryUrl), "application/json", reader)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	err = v.UpdateData(buf)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (v *Sentry) setBeam(beam Beam) error {
	err := v.requestBeam(beam)
	if err != nil {
		return err
	}

	return nil
}

func (v *Sentry) setPositions(position Position) error {
	pos := SetPosition{}

	pos.Tilt = position.Tilt - 90
	pos.Pan = position.Pan - 90
	pos.Token = os.Getenv("sentryToken")

	err := v.requestPosition(pos)
	if err != nil {
		return err
	}
	return nil

}

func (v *Sentry) listen() error {
	for {
		select {
		case beam := <-v.beamChannel:
			b := Beam{}
			err := json.Unmarshal([]byte(beam.Request), &b)
			if err != nil {
				log.Err(err)
			}
			err = v.setBeam(b)
			if err != nil {
				log.Err(err)
			}
		case position := <-v.positionChannel:
			p := Position{}
			err := json.Unmarshal([]byte(position.Request), &p)
			if err != nil {
				log.Err(err)
			}
			err = v.setPositions(p)
			if err != nil {
				log.Err(err)
			}
		}
	}
}

func (v *Sentry) pull() error {

	client := http.Client{}
	client.Timeout = time.Millisecond * 250
	get, err := client.Get(fmt.Sprintf("http://%s/status", sentryUrl))
	if err != nil {
		return nil
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(get.Body)
	if err != nil {
		return err
	}

	err = v.UpdateData(buf)
	if err != nil {
		return err
	}

	return nil
}

func (v *Sentry) UpdateData(buf bytes.Buffer) error {
	s := Status{}
	err := json.Unmarshal(buf.Bytes(), &s)
	if err != nil {
		return err
	}
	p := Position{}
	p.Pan = int(mapRange(float64(s.Servos.Pan), -90, 90, 0, 180))
	p.Tilt = int(mapRange(float64(s.Servos.Tilt), -90, 90, 0, 180))
	marshal, err := json.Marshal(p)
	if err != nil {
		return err
	}
	err = v.Attributes.Set(v.eId, "position", string(marshal))
	if err != nil {
		return err
	}

	b := Beam{}
	b.Target = "primary"
	if s.Beams.Primary.Active {
		b.Active = 1
	} else {
		b.Active = 0
	}

	b.Power = s.Beams.Primary.Power

	marshal, err = json.Marshal(b)
	if err != nil {
		return err
	}
	err = v.Attributes.Set(v.eId, "beam", string(marshal))
	if err != nil {
		return err
	}
	return nil
}

func (v *Sentry) Update() error {
	if v.Ready() {
		err := v.pull()
		if err != nil {
			v.Err(err)
		}
	}
	return nil
}

func (v *Sentry) Run() error {

	e := &domain.Entity{
		Name:   "sentryA",
		Module: "sentry",
		Type:   "media",
	}
	err := v.Entities.Register(e)
	if err != nil {
		return err
	}

	position := &domain.Attribute{
		Key:     "position",
		Value:   v.position.Marshal(),
		Request: v.position.Marshal(),
		Type:    "media",
		Order:   0,
		Entity:  e.Id,
		Channel: make(chan domain.Attribute, 10),
	}

	v.positionChannel = position.Channel

	v.beam = Beam{
		Target: "primary",
		Active: 0,
		Power:  15,
	}

	beam := &domain.Attribute{
		Key:     "beam",
		Value:   v.beam.Marshal(),
		Request: v.beam.Marshal(),
		Type:    "media",
		Order:   0,
		Entity:  e.Id,
		Channel: make(chan domain.Attribute),
	}

	v.beamChannel = beam.Channel

	v.eId = e.Id

	go func() {
		err = v.listen()
		if err != nil {
			return
		}
	}()

	err = v.Attributes.Register(position)
	if err != nil {
		return err
	}

	err = v.Attributes.Register(beam)
	if err != nil {
		return err
	}

	err = v.pull()
	if err != nil {
		return err
	}

	return nil
}
