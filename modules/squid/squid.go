// Copyright (c) 2021 Braden Nicholson

package main

import (
	"fmt"
	"strconv"
	"sync"
	"udap/internal/log"
	"udap/internal/models"
	"udap/internal/pkg/dmx"
	"udap/internal/pkg/dmx/ft232"
	"udap/pkg/plugin"
)

var Module Squid

type Squid struct {
	plugin.Module
	dmx        ft232.DMXController
	state      map[int]int
	stateMutex sync.RWMutex
	connected  bool
}

func init() {
	config := plugin.Config{
		Name:        "squid",
		Type:        "module",
		Description: "Control LOR Light Controller",
		Version:     "2.0.1",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (s *Squid) findDevices() error {
	if !s.connected {
		return nil
	}
	for i := 1; i <= 16; i++ {
		name := fmt.Sprintf("ch%d", i)
		entity := models.NewDimmer(name, s.Name)
		res, err := s.Entities.Register(entity)
		if err != nil {
			return err
		}

		on := models.Attribute{
			Key:     "on",
			Value:   "false",
			Request: "false",
			Type:    "toggle",
			Order:   0,
			Entity:  res.Id,
		}

		on.FnGet(func() (string, error) {
			state := "off"
			s.stateMutex.Lock()
			if s.state[i] > 0 {
				state = "on"
			}
			s.stateMutex.Unlock()
			return state, nil
		})

		on.FnPut(func(value string) error {
			s.stateMutex.Lock()
			if value == "on" {
				s.state[i] = 100
			} else {
				s.state[i] = 0
			}
			s.stateMutex.Unlock()
			return nil
		})

		err = s.Attributes.Register(&on)
		if err != nil {
			return err
		}

		dim := models.Attribute{
			Key:     "dim",
			Value:   "0",
			Request: "0",
			Type:    "range",
			Order:   1,
			Entity:  res.Id,
		}

		dim.FnGet(func() (string, error) {
			state := "0"

			out, err := s.dmx.GetChannel(int16(i))
			if err != nil {
			}
			state = fmt.Sprintf("%d", (out/255.0)*100)

			return state, nil
		})

		dim.FnPut(func(value string) error {
			s.stateMutex.Lock()
			parseInt, err := strconv.ParseInt(value, 10, 8)
			if err != nil {
				return err
			}
			s.state[i] = int((parseInt / 100.0) * 255.0)
			s.stateMutex.Unlock()
			return nil
		})

		err = s.Attributes.Register(&dim)
		if err != nil {
			return err
		}

		s.stateMutex.Lock()
		s.state[i] = 0
		s.stateMutex.Unlock()

	}

	return nil
}

// Setup is called once at the launch of the module
func (s *Squid) Setup() (plugin.Config, error) {
	return s.Config, nil
}
func (s *Squid) Setup2() (plugin.Config, error) {
	s.connected = false
	config := dmx.NewConfig(0x02)
	config.GetUSBContext()

	s.dmx = ft232.NewDMXController(config)

	defer func() {
		if r := recover(); r != nil {
			s.connected = false
			log.Err(fmt.Errorf("DMX Panicked: %s", r))
			return
		}
	}()
	err := s.dmx.Connect()
	if err != nil {
		return s.Config, err
	}
	s.stateMutex = sync.RWMutex{}
	s.state = map[int]int{}
	s.connected = true
	err = s.findDevices()
	if err != nil {
		return s.Config, err
	}

	return s.Config, nil
}

func (s *Squid) HandleState(channel int) func(state interface{}) error {
	return func(state interface{}) error {
		s.stateMutex.Lock()
		if state == nil {
			s.state[channel] = 0
		} else {

			parseInt, err := strconv.ParseInt(state.(string), 10, 16)
			if err != nil {
				return err
			}
			s.state[channel] = int(parseInt)
		}
		s.stateMutex.Unlock()
		return nil
	}
}

func (s *Squid) ReadState(channel int) func() interface{} {
	return func() interface{} {
		out, err := s.dmx.GetChannel(int16(channel))
		if err != nil {
			return ""
		}
		return out
	}
}

// Update is called every cycle
func (s *Squid) Update() error {

	return nil
}

// Run is called after Setup, concurrent with Update
func (s *Squid) Run() (err error) {

	return nil
}
