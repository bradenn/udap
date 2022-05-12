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
	entities   map[int]string
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

// setChannelValue sends a dmx signal to the provide channel with the provided value
func (s *Squid) setChannelValue(channel int, value int) (err error) {
	if !s.connected {
		return fmt.Errorf("squid is not connected")
	}

	if value > 100 || value < 0 {
		return fmt.Errorf("desired value '%d' is invalid", value)
	}

	adjustedValue := (value / 100.0) * 255

	err = s.dmx.SetChannel(int16(channel), byte(adjustedValue))
	if err != nil {
		return err
	}

	return nil
}

// getChannelValue polls the dmx controller for the current value of the channel
func (s *Squid) getChannelValue(channel int) (value int, err error) {
	if !s.connected {
		return 0, fmt.Errorf("squid is not connected")
	}

	res, err := s.dmx.GetChannel(int16(channel))
	if err != nil {
		return 0, err
	}
	newValue := (res / 255.0) * 100.0
	return int(newValue), nil
}

func (s *Squid) isLocalOn(channel int) (value bool) {
	value = false
	s.stateMutex.Lock()
	if s.state[channel] != 0 {
		value = true
	}
	s.stateMutex.Unlock()
	return value
}

func (s *Squid) getLocalValue(channel int) (value int) {
	value = 0
	s.stateMutex.Lock()
	value = s.state[channel]
	s.stateMutex.Unlock()
	return value
}

func (s *Squid) setLocalValue(channel int, value int) error {

	s.stateMutex.Lock()
	s.state[channel] = value
	s.stateMutex.Unlock()

	return nil
}

// isChannelOn provides a boolean describing the on state of the channel
func (s *Squid) isChannelOn(channel int) (value bool, err error) {
	if !s.connected {
		return false, fmt.Errorf("squid is not connected")
	}
	channelValue, err := s.getChannelValue(channel)
	if err != nil {
		return false, err
	}
	if channelValue > 0 {
		return true, nil
	}
	return false, nil
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

		s.entities[i] = res.Id

		on.FnGet(func() (string, error) {
			state := "off"
			if s.isLocalOn(i) {
				state = "on"
			}
			return state, nil
		})

		on.FnPut(func(value string) error {

			if value == "on" {
				err = s.setLocalValue(i, 100)
				if err != nil {
					return err
				}
			} else {
				err = s.setLocalValue(i, 0)
				if err != nil {
					return err
				}
			}

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
			value := s.getLocalValue(i)
			if err != nil {
				return "", err
			}
			state = fmt.Sprintf("%d", value)

			return state, nil
		})

		dim.FnPut(func(value string) error {

			parseInt, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return err
			}

			err = s.setLocalValue(i, int(parseInt))
			if err != nil {
				return err
			}

			return nil
		})

		err = s.Attributes.Register(&dim)
		if err != nil {
			return err
		}

	}

	return nil
}

// Setup is called once at the launch of the module
func (s *Squid) Setup() (plugin.Config, error) {

	return s.Config, nil
}

func (s *Squid) connect() error {
	s.connected = false
	config := dmx.NewConfig(0x02)
	config.GetUSBContext()
	defer func() {
		if r := recover(); r != nil {
			s.connected = false
			log.Err(fmt.Errorf("DMX Panicked: %s", r))
			return
		}
	}()
	s.dmx = ft232.NewDMXController(config)

	err := s.dmx.Connect()
	if err != nil {
		return err
	}

	s.stateMutex = sync.RWMutex{}
	s.stateMutex.Lock()
	s.state = map[int]int{}
	s.stateMutex.Unlock()
	s.entities = map[int]string{}
	s.connected = true
	err = s.findDevices()
	if err != nil {
		return err
	}
	return nil
}

func (s *Squid) pull() error {
	if !s.connected {
		return nil
	}

	for i, entity := range s.entities {
		state := "off"

		if s.isLocalOn(i) {
			state = "on"
		}
		err := s.Attributes.Set(entity, "on", state)
		if err != nil {
			return err
		}

		state = "0"
		value := s.getLocalValue(i)
		state = fmt.Sprintf("%d", value)
		err = s.Attributes.Set(entity, "dim", state)
		if err != nil {
			return err
		}
	}

	return nil

}

// Update is called every cycle
func (s *Squid) Update() error {
	// pulse.Fixed(1000)
	// defer pulse.End()
	// if time.Since(s.Module.LastUpdate) >= time.Second*10 {
	// 	s.Module.LastUpdate = time.Now()
	// 	return s.pull()
	// }
	return nil
}

// Run is called after Setup, concurrent with Update
func (s *Squid) Run() (err error) {
	//
	// err = s.connect()
	// if err != nil {
	// 	return err
	// }
	//
	// for {
	// 	s.stateMutex.RLock()
	// 	for k, v := range s.state {
	// 		err = s.dmx.SetChannel(int16(k), byte(v))
	// 		if err != nil {
	// 			log.Err(err)
	// 		}
	// 	}
	// 	s.stateMutex.RUnlock()
	//
	// 	err = s.dmx.Render()
	// 	if err != nil {
	// 		log.Err(err)
	// 	}
	// 	time.Sleep(20 * time.Millisecond)
	// }
	return nil
}
