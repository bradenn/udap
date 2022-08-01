// Copyright (c) 2021 Braden Nicholson

package main

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
	"udap/platform/dmx"
)

var Module Squid

type Squid struct {
	plugin.Module
	dmx      dmx.DMXController
	state    map[int]int
	entities map[int]string
	receiver chan domain.Attribute
	mutex    sync.RWMutex

	update chan Command

	connected bool
}

type Command struct {
	read  bool
	id    int
	value int
}

func init() {
	config := plugin.Config{
		Name:        "squid",
		Type:        "module",
		Description: "Control LOR Light Controller",
		Version:     "2.1",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (s *Squid) setChannelValue(channel int, value int) (err error) {
	if !s.connected {
		return fmt.Errorf("squid is not connected")
	}

	if value > 100 || value < 0 {
		return fmt.Errorf("desired value '%d' is invalid", value)
	}
	var adjustedValue byte
	adjustedValue = uint8(math.Round((float64(value) / 100.0) * 255.0))

	err = s.dmx.SetChannel(int16(channel), adjustedValue)
	if err != nil {
		return err
	}

	return nil
}

//
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

func (s *Squid) getLocalValue(channel int) (value int) {
	level := 0
	s.mutex.RLock()
	level = s.state[channel]
	s.mutex.RUnlock()

	return level
}

func (s *Squid) setLocalValue(channel int, value int) error {

	payload := Command{
		read:  false,
		id:    channel,
		value: value,
	}

	select {
	case s.update <- payload:
		return nil
	case <-time.After(time.Millisecond * 100):
		s.WarnF("setting local value for channel %d timed out", channel)
	}

	return nil
}

func (s *Squid) handleDimState(channel int, value string) error {
	parseInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	s.setValue(channel, int(parseInt))
	return nil
}

func (s *Squid) setValue(channel int, value int) {
	s.LogF("Setting %d to %d", channel, value)
	s.mutex.Lock()
	s.state[channel] = value
	s.mutex.Unlock()
}

func (s *Squid) handleOnState(channel int, value string) error {
	if value == "true" {
		s.setValue(channel, 100)
	} else {
		s.setValue(channel, 0)
	}
	return nil
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
	if !s.connected {
		return nil
	}
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
		err = s.setLocalValue(i, 0)
		if err != nil {
			return err
		}

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
	return s.Config, nil
}

func (s *Squid) connect() error {
	s.connected = false

	config := dmx.NewConfig(0x02)
	config.GetUSBContext()

	defer func() {
		if r := recover(); r != nil {
			s.connected = false
			s.ErrF("DMX Panicked: %s", r)
			return
		}
	}()

	s.dmx = dmx.NewDMXController(config)

	err := s.dmx.Connect()
	if err != nil {
		return err
	}

	s.connected = true
	s.update = make(chan Command)

	s.receiver = make(chan domain.Attribute)

	s.state = map[int]int{}

	s.entities = map[int]string{}

	s.mutex = sync.RWMutex{}

	go s.mux()

	err = s.registerDevices()
	if err != nil {
		return err
	}

	return nil
}

func (s *Squid) mux() {
	for {
		select {
		case cmd := <-s.update:
			s.setValue(cmd.id, cmd.value)
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
		s.WarnF("not connected")
		return nil
	}

	for i, entity := range s.entities {

		stateNum := "0"
		value := s.getLocalValue(i)
		stateNum = fmt.Sprintf("%d", value)
		err := s.Attributes.Update(entity, "dim", stateNum, time.Now())
		if err != nil {
			return err
		}

		state := "false"
		if value > 0 {
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

func (s *Squid) renderer() {
	for {

		if !s.connected {
			s.ErrF("not connected, exiting render loop")
			break
		}
		time.Sleep(50 * time.Millisecond)
		for i := 1; i <= 16; i++ {
			s.mutex.RLock()
			err := s.setChannelValue(i, s.state[i])
			s.mutex.RUnlock()
			if err != nil {
				s.ErrF("failed to set the value for channel %d: %s", i, err.Error())
			}
		}

		err := s.dmx.Render()
		if err != nil {
			s.ErrF("%s", err.Error())
			break
		}

	}

	err := s.dmx.Close()
	if err != nil {
		s.ErrF("%s", err.Error())
	}
}

// Run is called after Setup, concurrent with Update
func (s *Squid) Run() (err error) {

	err = s.connect()
	if err != nil {
		return err
	}

	go s.renderer()

	return nil
}
