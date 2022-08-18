// Copyright (c) 2021 Braden Nicholson

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
)

var Module Squid

type Channel struct {
	Channel int `json:"channel"`
	Value   int `json:"value"`
}

func (c *Channel) isOn() bool {
	return c.Value > 0
}

func (c *Channel) toPercent() int {
	return int((float64(c.Value) / 255.0) * 100.0)
}

func (c *Channel) fromPercent(percent int) int {
	return int((float64(percent) / 100.0) * 255.0)
}

type Status struct {
	Channels []Channel `json:"channels"`
}

type Squid struct {
	plugin.Module
	entities map[int]string
	receiver chan domain.Attribute

	connected bool
}

func init() {
	config := plugin.Config{
		Name:        "squid",
		Type:        "module",
		Description: "Control LOR Light Controller",
		Version:     "3.0",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (s *Squid) setChannelValue(channel int, value int) (err error) {
	if value > 100 || value < 0 {
		return fmt.Errorf("desired value '%d' is invalid", value)
	}
	var adjustedValue int
	adjustedValue = int(math.Round((float64(value) / 100.0) * 255.0))

	err = s.remoteRequest(channel, adjustedValue)
	if err != nil {
		return err
	}

	return nil
}

func (s *Squid) handleDimState(channel int, value string) error {
	parseInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}

	err = s.setChannelValue(channel, int(parseInt))
	if err != nil {
		return err
	}

	return nil
}

func (s *Squid) handleOnState(channel int, value string) error {
	if value == "true" {
		err := s.setChannelValue(channel, 100)
		if err != nil {
			return err
		}
	} else {
		err := s.setChannelValue(channel, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Squid) remoteRequest(channel int, value int) error {
	c := Channel{
		Channel: channel,
		Value:   value,
	}
	marshal, err := json.Marshal(c)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.Write(marshal)
	client := http.Client{Timeout: time.Millisecond * 500}
	post, err := client.Post("http://10.0.1.76/channel", "application/json", &buf)

	defer func() {
		post.Body.Close()
	}()

	if err != nil {
		return err
	}

	var response bytes.Buffer
	_, err = response.ReadFrom(post.Body)
	if err != nil {
		return err
	}

	if response.String() != "okay" {
		return fmt.Errorf("request failed")
	}

	return nil

}

func (s *Squid) remoteFetch() (Status, error) {
	client := http.Client{Timeout: time.Millisecond * 500}
	get, err := client.Get("http://10.0.1.76/status")
	defer get.Body.Close()
	if err != nil {
		return Status{}, err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(get.Body)
	if err != nil {
		return Status{}, err
	}

	status := Status{}
	err = json.Unmarshal(buf.Bytes(), &status)
	if err != nil {
		return Status{}, err
	}

	return status, err
}

func (s *Squid) handleAttribute(attr domain.Attribute) error {
	channel := 0
	for i, id := range s.entities {
		if attr.Entity == id {
			channel = i
		}
	}
	if channel <= 0 {
		return nil
	}
	if attr.Key == "on" {
		err := s.handleOnState(channel, attr.Request)
		if err != nil {
			return err
		}
	} else if attr.Key == "dim" {
		err := s.handleDimState(channel, attr.Request)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Squid) registerDevices() error {
	s.entities = map[int]string{}
	s.receiver = make(chan domain.Attribute)
	for i := 1; i <= 16; i++ {

		entity := &domain.Entity{
			Name:   fmt.Sprintf("ch%d", i),
			Module: s.Config.Name,
			Type:   "dimmer",
		}
		err := s.Entities.Register(entity)
		if err != nil {
			return err
		}

		on := &domain.Attribute{
			Key:     "on",
			Value:   "false",
			Request: "false",
			Type:    "toggle",
			Order:   0,
			Entity:  entity.Id,
			Channel: s.receiver,
		}

		s.entities[i] = entity.Id

		err = s.Attributes.Register(on)
		if err != nil {
			return err
		}

		dim := &domain.Attribute{
			Key:     "dim",
			Value:   "0",
			Request: "0",
			Type:    "range",
			Order:   1,
			Entity:  entity.Id,
			Channel: s.receiver,
		}

		err = s.Attributes.Register(dim)
		if err != nil {
			return err
		}

	}

	return nil
}

// Setup is called once at the launch of the module
func (s *Squid) Setup() (plugin.Config, error) {
	err := s.UpdateInterval(1000)
	if err != nil {
		return plugin.Config{}, err
	}
	err = s.InitConfig("address", "10.0.1.76")
	if err != nil {
		return plugin.Config{}, err
	}
	s.connected = false
	return s.Config, nil
}

func (s *Squid) mux() {
	for {
		select {
		case attr := <-s.receiver:
			err := s.handleAttribute(attr)
			if err != nil {
				return
			}
		}

	}
}

func (s *Squid) pull() error {
	if !s.connected {
		return nil
	}
	fetch, err := s.remoteFetch()
	if err != nil {
		return err
	}

	for _, channel := range fetch.Channels {
		entity := s.entities[channel.Channel]

		dimValue := fmt.Sprintf("%d", channel.toPercent())
		err = s.Attributes.Update(entity, "dim", dimValue, time.Now())
		if err != nil {
			return err
		}

		state := "false"
		if channel.isOn() {
			state = "true"
		}
		err = s.Attributes.Update(entity, "on", state, time.Now())
		if err != nil {
			return err
		}
	}

	return nil

}

// Update is called every cycle
func (s *Squid) Update() error {
	if s.Ready() {
		return s.pull()
	}
	return nil
}

// Run is called after Setup, concurrent with Update
func (s *Squid) Run() (err error) {
	err = s.registerDevices()
	if err != nil {
		return err
	}
	go s.mux()
	s.connected = true
	return nil
}
