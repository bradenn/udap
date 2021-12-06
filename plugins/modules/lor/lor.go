// Copyright (c) 2021 Braden Nicholson

package main

import (
	"context"
	"fmt"
	"strconv"
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
}

func init() {
	config := plugin.Config{
		Name:        "squid",
		Type:        "module",
		Description: "Lor 16-Channel 120VAC 15A/15A 10-Bit Dimmer",
		Version:     "0.1.1",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (s *Lor) Setup() (plugin.Config, error) {
	config := dmx.NewConfig(0x02)

	config.GetUSBContext()
	// Since ft232 is a shitty module, it panics when USB can't be found.
	defer func() {
		recover()
	}()

	s.dmx = ft232.NewDMXController(config)
	if err := s.dmx.Connect(); err != nil {
		fmt.Printf("failed to connect DMX Controller: %s\n", err)
	}
	s.stateMutex = sync.RWMutex{}
	s.state = map[int]int{}
	//
	// go func() {
	// 	defer func() {
	// 		err := controller.Close()
	// 		if err != nil {
	// 			return
	// 		}
	// 	}()
	// 	var sequence []rune
	//
	// 	morseStr := "I LOVE YOU"
	// 	// delay := 0
	// 	log.Sherlock("Running Morse: '%s'", morseStr)
	// 	// For each letter of the string
	// 	for i := range morseStr {
	// 		morse := morseSymbols[string(morseStr[i])]
	// 		for _, symbol := range morse {
	// 			sequence = append(sequence, symbol)
	// 		}
	// 	}
	//
	// 	fmt.Println(sequence)
	//
	// 	// place := 0
	// 	for {
	//
	// 		var err error
	//
	// 		err = controller.SetChannel(1, byte(255))
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	//
	// 		err = controller.Render()
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	//
	// 		time.Sleep(time.Millisecond * 250)
	// 		// _ = controller.Render()
	// 		// switch r := sequence[place]; r {
	// 		// case '.':
	// 		// 	time.Sleep(time.Millisecond * 250)
	// 		// 	break
	// 		// case '-':
	// 		// 	time.Sleep(time.Millisecond * 500)
	// 		// 	break
	// 		// case ' ':
	// 		// 	time.Sleep(time.Millisecond * 250)
	// 		// 	break
	// 		// }
	//
	// 		// time.Sleep(time.Millisecond * 100)
	//
	// 	}
	// }()

	return s.Config, nil
}

func (s *Lor) Run() error {

	defer func() {
		err := s.dmx.Close()
		if err != nil {
			return
		}
	}()

	for i := 1; i <= 16; i++ {
		name := fmt.Sprintf("ch%d", i)
		dimmer := models.NewDimmer(name, s.Name)
		dimmer.SetStateHandler(s.HandleState(i))
		dimmer.SetStateReceiver(s.ReadState(i))
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

	return nil
}

func (s *Lor) HandleState(channel int) func(state interface{}) error {
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

func (s *Lor) ReadState(channel int) func() interface{} {
	return func() interface{} {
		out, err := s.dmx.GetChannel(int16(channel))
		if err != nil {
			return ""
		}
		return out
	}
}

func (s *Lor) Update(ctx context.Context) error {
	return nil
}
