// Copyright (c) 2021 Braden Nicholson

package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"udap/internal/log"
	"udap/internal/models"
	"udap/internal/pkg/dmx"
	"udap/internal/pkg/dmx/ft232"
	"udap/pkg/plugin"
)

var Module Lor

type Lor struct {
	plugin.Module
	dmx        ft232.DMXController
	dmxConf    dmx.ControllerConfig
	state      map[int]int
	stateMutex sync.Mutex
	open       bool
}

func init() {
	config := plugin.Config{
		Name:        "squid",
		Type:        "module",
		Description: "Lor 16-Channel 120VAC 15A/15A 10-Bit Dimmer",
		Version:     "0.1.1",
		Author:      "Braden Nicholson",
	}
	Module.open = false
	Module.Config = config
}

func (s *Lor) Setup() (plugin.Config, error) {

	s.state = map[int]int{}
	s.stateMutex = sync.Mutex{}

	s.stateMutex.Lock()
	s.dmxConf = dmx.NewConfig(0x02)
	s.dmxConf.GetUSBContext()
	s.stateMutex.Unlock()
	// Since ft232 is a shitty module, it panics when USB can't be found.

	return s.Config, nil
}

func (s *Lor) Update() error {
	return nil
}

func (s *Lor) Run() error {
	defer func() {
		recover()
	}()
	defer func() {
		if s.open {
			err := s.dmx.Close()
			if err != nil {
				return
			}
		}
	}()

	s.dmx = ft232.NewDMXController(s.dmxConf)
	if err := s.dmx.Connect(); err != nil {
		fmt.Printf("failed to connect DMX Controller: %s\n", err)
	} else {
		Module.open = true
	}

	if !s.open {
		s.dmx = ft232.DMXController{}
		return fmt.Errorf("squid is not connected")
	}

	for i := 1; i <= 16; i++ {
		name := fmt.Sprintf("ch%d", i)
		dimmer := models.NewDimmer(name, s.Name)
		err := dimmer.Handlers(s.Tx(i), s.Rx(i))
		if err != nil {
			continue
		}
		s.RegisterEntity(dimmer)
	}

	for {
		s.stateMutex.Lock()
		for id, value := range s.state {
			err := s.dmx.SetChannel(int16(id), byte(value))
			if err != nil {
				log.Err(err)
			}
		}
		s.stateMutex.Unlock()
		err := s.dmx.Render()
		if err != nil {
			log.Err(err)
		}

		time.Sleep(time.Millisecond * 25)
	}
}

func (s *Lor) Tx(channel int) models.Tx {
	return func(state models.State) error {
		mono := models.Mono{}
		err := json.Unmarshal(state, &mono)
		if err != nil {
			return err
		}
		s.stateMutex.Lock()
		s.state[channel] = int(mono.Value * 255.0)
		s.stateMutex.Unlock()
		return nil
	}
}

func (s *Lor) Rx(channel int) models.Rx {
	return func() models.State {
		mono := models.Mono{}
		s.stateMutex.Lock()
		mono.Value = float32(s.state[channel]) / 255.0
		s.stateMutex.Unlock()
		return mono.Marshal()
	}
}
