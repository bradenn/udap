// Copyright (c) 2021 Braden Nicholson

package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
)

var Module Squid

type Squid struct {
	plugin.Module
	// dmx        ft232.DMXController
	state      map[int]int
	entities   map[int]string
	stateMutex sync.Mutex
	update     chan Command

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
		Version:     "1.2",
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
	// var adjustedValue byte
	// adjustedValue = uint8(math.Round((float64(value) / 100.0) * 255.0))
	//
	// err = s.dmx.SetChannel(int16(channel), adjustedValue)
	// if err != nil {
	// 	return err
	// }

	return nil
}

//
// // getChannelValue polls the dmx controller for the current value of the channel
// func (s *Squid) getChannelValue(channel int) (value int, err error) {
// 	if !s.connected {
// 		return 0, fmt.Errorf("squid is not connected")
// 	}
//
// 	res, err := s.dmx.GetChannel(int16(channel))
// 	if err != nil {
// 		return 0, err
// 	}
// 	newValue :=
// 	return int(newValue), nil
// }

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

	s.update <- Command{
		read:  false,
		id:    channel,
		value: value,
	}

	return nil
}

func (s *Squid) remoteGetOn(id int) func() (string, error) {
	return func() (string, error) {
		state := "false"
		if s.isLocalOn(id) {
			state = "true"
		}
		return state, nil
	}
}

func (s *Squid) remotePutOn(id int) func(value string) error {
	return func(value string) error {
		if value == "true" {
			err := s.setLocalValue(id, 100)
			if err != nil {
				return err
			}
		} else {
			err := s.setLocalValue(id, 0)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (s *Squid) remoteGetDim(id int) func() (string, error) {
	return func() (string, error) {
		state := "0"
		value := s.getLocalValue(id)
		state = fmt.Sprintf("%d", value)

		return state, nil
	}
}

func (s *Squid) remotePutDim(id int) func(value string) error {
	return func(value string) error {

		parseInt, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return err
		}

		err = s.setLocalValue(id, int(parseInt))
		if err != nil {
			return err
		}

		return nil
	}
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

		on := domain.Attribute{
			Key:     "on",
			Value:   "false",
			Request: "false",
			Type:    "toggle",
			Order:   0,
			Entity:  entity.Id,
			Channel: make(chan domain.Attribute),
		}

		s.entities[i] = entity.Id

		go func() {
			for attribute := range on.Channel {
				err = s.remotePutOn(i)(attribute.Request)
				if err != nil {
					return
				}
			}
		}()

		err = s.Attributes.Register(&on)
		if err != nil {
			return err
		}

		dim := domain.Attribute{
			Key:     "dim",
			Value:   "0",
			Request: "0",
			Type:    "range",
			Order:   1,
			Entity:  entity.Id,
			Channel: make(chan domain.Attribute),
		}

		go func() {
			for attribute := range dim.Channel {
				err = s.remotePutDim(i)(attribute.Request)
				if err != nil {
					return
				}
			}
		}()

		err = s.Attributes.Register(&dim)
		if err != nil {
			return err
		}

	}

	return nil
}

// Setup is called once at the launch of the module
func (s *Squid) Setup() (plugin.Config, error) {
	err := s.UpdateInterval(2000)
	if err != nil {
		return plugin.Config{}, err
	}
	return s.Config, nil
}

// func (s *Squid) connect() error {
// 	s.connected = false
//
// 	so := os.Stdout
// 	config := dmx.NewConfig(0x02)
// 	config.GetUSBContext()
//
// 	os.Stdout = so
//
// 	defer func() {
// 		if r := recover(); r != nil {
// 			s.connected = false
// 			log.Err(fmt.Errorf("DMX Panicked: %s", r))
// 			return
// 		}
// 	}()
//
// 	s.dmx = ft232.NewDMXController(config)
//
// 	err := s.dmx.Connect()
// 	if err != nil {
// 		return err
// 	}
//
// 	s.connected = true
// 	s.update = make(chan Command, 4)
//
// 	s.stateMutex = sync.Mutex{}
// 	s.stateMutex.Lock()
// 	s.state = map[int]int{}
// 	s.stateMutex.Unlock()
//
// 	s.entities = map[int]string{}
//
// 	go s.mux()
//
// 	err = s.registerDevices()
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func (s *Squid) mux() {
	for cmd := range s.update {
		s.stateMutex.Lock()
		s.state[cmd.id] = cmd.value
		s.stateMutex.Unlock()
	}
}

func (s *Squid) pull() error {
	if !s.connected {
		return nil
	}

	for i, entity := range s.entities {
		state := "false"

		if s.isLocalOn(i) {
			state = "true"
		}
		err := s.Attributes.Update(entity, "on", state, time.Now())
		if err != nil {
			return err
		}

		state = "0"
		value := s.getLocalValue(i)
		state = fmt.Sprintf("%d", value)
		err = s.Attributes.Update(entity, "dim", state, time.Now())
		if err != nil {
			return err
		}
	}

	return nil

}

// Update is called every cycle
func (s *Squid) Update() error {
	if s.Ready() {
		err := s.UpdateInterval(2000)
		if err != nil {
			return err
		}
		err = s.pull()
		if err != nil {
			return err
		}
	}
	return nil
}

// Run is called after Setup, concurrent with Update
func (s *Squid) Run() (err error) {

	// err = s.connect()
	// if err != nil {
	// 	return err
	// }
	//
	// for {
	//
	// 	if !s.connected {
	// 		break
	// 	}
	//
	// 	s.stateMutex.Lock()
	// 	for k, v := range s.state {
	// 		err = s.setChannelValue(k, v)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	s.stateMutex.Unlock()
	//
	// 	err = s.dmx.Render()
	// 	if err != nil {
	// 		log.Err(err)
	// 		break
	// 	}
	// 	time.Sleep(50 * time.Millisecond)
	//
	// }
	//
	// err = s.dmx.Close()
	// if err != nil {
	// 	return err
	// }

	return nil
}
