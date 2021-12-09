// Copyright (c) 2021 Braden Nicholson

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
	"udap/internal/models"
	"udap/internal/pkg/dmx"
	"udap/internal/pkg/dmx/ft232"
	"udap/pkg/plugin"
)

var Module Lor

type Lor struct {
	plugin.Module
	dmx        ft232.DMXController
	state      map[int]int
	stateMutex sync.RWMutex
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
	config := dmx.NewConfig(0x02)
	s.state = map[int]int{}
	config.GetUSBContext()
	// Since ft232 is a shitty module, it panics when USB can't be found.
	defer func() {
		recover()
	}()

	s.dmx = ft232.NewDMXController(config)
	if err := s.dmx.Connect(); err != nil {
		fmt.Printf("failed to connect DMX Controller: %s\n", err)
	} else {
		Module.open = true
	}
	s.stateMutex = sync.RWMutex{}
	s.state = map[int]int{}

	return s.Config, nil
}

func (s *Lor) Run() error {

	if !s.open {
		s.dmx = ft232.DMXController{}
		return fmt.Errorf("squid is not connected")
	}

	defer func() {
		err := s.dmx.Close()
		if err != nil {
			return
		}

	}()

	for i := 1; i <= 16; i++ {
		name := fmt.Sprintf("ch%d", i)
		dimmer := models.NewDimmer(name, s.Name)
		err := dimmer.Handlers(s.HandleState(i), s.ReadState(i))
		if err != nil {
			return err
		}
		s.stateMutex.Lock()
		s.state[i] = 0
		s.stateMutex.Unlock()
		s.RegisterEntity(dimmer)
	}

	for {
		for id, value := range s.state {
			err := s.dmx.SetChannel(int16(id), byte(value))
			if err != nil {
				return err
			}
		}
		err := s.dmx.Render()
		if err != nil {
			return err
		}

		time.Sleep(time.Millisecond * 20)
	}
}

func (s *Lor) HandleState(channel int) func(state models.State) error {
	return func(state models.State) error {
		s.stateMutex.Lock()
		if state == nil {
			s.state[channel] = 0
		} else {
			mono := models.Mono{}
			mono.Unmarshal(state)
			out := int(mono.Value * 255)
			s.state[channel] = out
		}
		s.stateMutex.Unlock()
		return nil
	}
}

func (s *Lor) ReadState(channel int) func() models.State {
	return func() models.State {
		mono := models.Mono{}
		out, err := s.dmx.GetChannel(int16(channel))
		if err != nil {
			return []byte{out}
		}
		mono.Value = float64(out / 255)
		return mono.Marshal()
	}
}

func (s *Lor) Update(ctx context.Context) error {
	return nil
}
