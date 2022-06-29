// Copyright (c) 2021 Braden Nicholson

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"udap/internal/core/domain"
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
}

type Beam struct {
	Target string `json:"target"`
	Active int    `json:"active"`
	Power  int    `json:"power"`
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

const sentryUrl = "10.0.1.60"

func init() {
	config := plugin.Config{
		Name:        "sentry",
		Type:        "module",
		Description: "Sentry with LASERS",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
	Module.eId = ""
	Module.Config = config
}

func (v *Sentry) Setup() (plugin.Config, error) {
	err := v.UpdateInterval(10000)
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
		Primary   bool `json:"primary"`
		Secondary bool `json:"secondary"`
	} `json:"beams"`
}

type SetPosition struct {
	Target   string `json:"target"`
	Position int    `json:"position"`
}

func mapRange(value float64, low1 float64, high1 float64, low2 float64, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

func (v *Sentry) requestPosition(position SetPosition) error {
	marshal, err := json.Marshal(position)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(marshal)

	client := http.Client{}
	_, err = client.Post(fmt.Sprintf("http://%s/position", sentryUrl), "application/json", reader)
	if err != nil {
		return err
	}

	return nil
}

func (v *Sentry) setPositions(position Position) error {
	pos := SetPosition{}
	pos.Target = "tilt"
	value := mapRange(float64(position.Pan), 0, 180, -256, 256)
	pos.Position = int(value)

	err := v.requestPosition(pos)
	if err != nil {
		return err
	}

	pos.Target = "pan"
	value = mapRange(float64(position.Tilt), 0, 180, -256, 256)
	pos.Position = int(value)

	err = v.requestPosition(pos)
	if err != nil {
		return err
	}

	return nil

}

func (v *Sentry) listen() error {
	for {
		select {
		case _ = <-v.beamChannel:

		case position := <-v.positionChannel:
			p := Position{}
			err := json.Unmarshal([]byte(position.Request), &p)
			if err != nil {
				return err
			}
			err = v.setPositions(p)
			if err != nil {
				return err
			}
		}
	}
}

func (v *Sentry) pull() error {

	client := http.Client{}
	get, err := client.Get(fmt.Sprintf("http://%s/status", sentryUrl))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(get.Body)
	if err != nil {
		return err
	}
	s := Status{}
	err = json.Unmarshal(buf.Bytes(), &s)
	if err != nil {
		return err
	}
	p := Position{}
	p.Pan = int(mapRange(float64(s.Servos.Tilt), -256, 256, 0, 180))
	p.Tilt = int(mapRange(float64(s.Servos.Pan), -256, 256, 0, 180))
	marshal, err := json.Marshal(p)
	if err != nil {
		return err
	}
	err = v.Attributes.Set(v.eId, "position", string(marshal))
	if err != nil {
		return err
	}

	return nil
}

func (v *Sentry) Update() error {
	if v.Ready() {
		err := v.pull()
		if err != nil {
			return err
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
		Channel: make(chan domain.Attribute),
	}

	v.positionChannel = position.Channel

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
