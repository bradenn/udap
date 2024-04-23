// Copyright (c) 2023 Braden Nicholson

package main

import (
	"fmt"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"math"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/plugin"
)

var Module Midi

type Midi struct {
	plugin.Module
	request chan domain.Attribute
	done    chan bool
}

type entry struct {
	Name string    `json:"name"`
	Time time.Time `json:"time"`
	Id   string    `json:"id"`
}

func init() {

	config := plugin.Config{
		Name:        "midi",
		Type:        "module",
		Description: "A midi module...",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}

	Module.Config = config
}

func (c *Midi) mux() {
	for attribute := range c.request {
		c.handleRequest(attribute)
	}
}

func (c *Midi) handleRequest(attribute domain.Attribute) {

}

func (c *Midi) Setup() (plugin.Config, error) {
	c.done = make(chan bool)
	return Module.Config, nil
}

func (c *Midi) Update() error {
	return nil
}

func (c *Midi) Listen() error {

	defer midi.CloseDriver()
	trigger := domain.Trigger{
		Name:        "midi-trigger",
		Type:        "module",
		Description: "Just a midi trigger",
	}
	err := c.Triggers.Register(&trigger)
	if err != nil {
		log.Err(err)
	}
	fmt.Printf("outports:%s\n ", midi.GetOutPorts())
	in, err := midi.FindInPort("MIDI-ME CASIO USB-MIDI")
	if err != nil {
		return err
	}
	for {
		stop, err := midi.ListenTo(in, func(msg midi.Message, timestampms int32) {
			var bt []byte
			var ch, key, vel uint8
			switch {
			case msg.GetSysEx(&bt):
				fmt.Printf("got sysex: % X\n", bt)
			case msg.GetNoteStart(&ch, &key, &vel):
				//err = c.Triggers.Trigger(&trigger)
				//if err != nil {
				//	log.Err(err)
				//}
				zones, err := c.Zones.FindByName("Office")
				if err != nil {
					return
				} // 68 132
				normalKey := float64(key-36.0) / (96.0 - 36.0)
				normalVel := float64(vel) / 128.0
				dim := int(math.Round(normalVel * 100.0))
				hue := int(math.Round(normalKey * 360.0))
				for _, entity := range zones.Entities {
					_ = c.Attributes.Request(entity.Id, "dim", fmt.Sprintf("%d", dim))
					_ = c.Attributes.Request(entity.Id, "hue", fmt.Sprintf("%d", hue))
				}
				fmt.Printf("dim %d, hue %d\n", dim, hue)
			default:
				// ignore
			}
		}, midi.UseSysEx())
		if err != nil {
			log.Err(err)
		}
		select {
		case <-c.done:
			stop()
		}
		time.Sleep(5000)
	}

	return nil
}

func (c *Midi) Run() error {
	c.request = make(chan domain.Attribute)

	go func() {
		err := c.Listen()
		if err != nil {
			log.Err(err)
			return
		}
	}()

	return nil
}

func (c *Midi) Dispose() error {
	close(c.request)
	c.done <- true
	return nil
}
