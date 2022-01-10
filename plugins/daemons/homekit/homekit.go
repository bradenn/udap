// Copyright (c) 2021 Braden Nicholson

package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/gerow/go-color"
	"log"
	log2 "udap/internal/log"
	"udap/internal/models"
	"udap/pkg/plugin"
)

var Module Homekit

type Homekit struct {
	plugin.Module
	Config plugin.Config
}

func (h *Homekit) Set(id string, mode string, value any) error {
	remote, err := h.SendId("entity", "find", id, nil)
	if err != nil {
		return err
	}

	entity := remote.(*models.Entity)
	state := models.LightState{}
	state.Parse(entity.State)
	rgbo := color.RGB{
		R: float64(state.Red),
		G: float64(state.Green),
		B: float64(state.Blue),
	}
	hslo := rgbo.ToHSL()

	switch mode {
	case "brightness":
		state.Level = value.(int)
		state.Mode = mode
		break
	case "power":
		power := value.(bool)
		if power {
			state.Power = "on"
		} else {
			state.Power = "off"
		}
		state.Mode = mode
		break
	case "hue":
		hsl := color.HSL{
			H: value.(float64),
			S: hslo.S,
			L: 50,
		}
		rgb := hsl.ToRGB()
		state.Red = int(rgb.R)
		state.Green = int(rgb.G)
		state.Blue = int(rgb.B)
		state.Mode = "color"
		break
	case "saturation":
		hsl := color.HSL{
			H: hslo.H,
			S: value.(float64),
			L: 50,
		}
		rgb := hsl.ToRGB()
		state.Red = int(rgb.R)
		state.Green = int(rgb.G)
		state.Blue = int(rgb.B)
		state.Mode = "color"
		break
	}
	entity.State = state.JSON()
	_, err = h.Send("entity", "state", entity.State)
	if err != nil {
		return err
	}
	return nil
}

func (h *Homekit) Get(id string, mode string) (any, error) {
	remote, err := h.SendId("entity", "find", id, nil)
	if err != nil {
		return nil, err
	}
	entity := remote.(*models.Entity)
	state := models.LightState{}
	state.Parse(entity.State)
	rgb := color.RGB{
		R: float64(state.Red),
		G: float64(state.Green),
		B: float64(state.Blue),
	}
	hsl := rgb.ToHSL()
	switch mode {
	case "brightness":
		return state.Level, nil
	case "power":
		return state.Power, nil
	case "hue":
		return hsl.H, nil
	case "saturation":
		return hsl.S, nil
	}

	return nil, nil
}

func (h *Homekit) Setup() (plugin.Config, error) {
	send, err := h.Send("entity", "compile", nil)
	if err != nil {
		return plugin.Config{}, err
	}

	entities := send.([]*models.Entity)

	var accessories []*accessory.Accessory

	for _, entity := range entities {
		if entity.Type != "wildcard" {
			continue
		}
		info := accessory.Info{
			Name:             entity.Name,
			ID:               uint64(entity.CreatedAt.UnixMilli()),
			SerialNumber:     entity.Id,
			Manufacturer:     entity.Module,
			Model:            entity.Type,
			FirmwareRevision: "1.0.1",
		}
		ac := accessory.NewColoredLightbulb(info)

		ac.Lightbulb.Brightness.OnValueRemoteUpdate(func(i int) {
			err = h.Set(entity.Id, "brightness", i)
			if err != nil {
				log2.Err(err)
				return
			}
		})

		ac.Lightbulb.Brightness.OnValueRemoteGet(func() int {
			i, err := h.Get(entity.Id, "brightness")
			if err != nil {
				log2.Err(err)
				return 0
			}
			return i.(int)
		})

		ac.Lightbulb.Hue.OnValueRemoteUpdate(func(i float64) {
			err = h.Set(entity.Id, "hue", i)
			if err != nil {
				log2.Err(err)
				return
			}
		})

		ac.Lightbulb.Hue.OnValueRemoteGet(func() float64 {
			i, err := h.Get(entity.Id, "hue")
			if err != nil {
				log2.Err(err)
				return 0
			}
			return i.(float64)
		})

		ac.Lightbulb.Saturation.OnValueRemoteUpdate(func(i float64) {
			err = h.Set(entity.Id, "saturation", i)
			if err != nil {
				log2.Err(err)
				return
			}
		})

		ac.Lightbulb.Saturation.OnValueRemoteGet(func() float64 {
			i, err := h.Get(entity.Id, "saturation")
			if err != nil {
				log2.Err(err)
				return 0
			}
			return i.(float64)
		})

		ac.Lightbulb.On.OnValueRemoteGet(func() bool {
			i, err := h.Get(entity.Id, "power")
			if err != nil {
				log2.Err(err)
				return false
			}
			return i.(string) == "on"
		})

		ac.Lightbulb.On.OnValueRemoteUpdate(func(b bool) {
			err = h.Set(entity.Id, "power", b)
			if err != nil {
				log2.Err(err)
				return
			}
		})

		accessories = append(accessories, ac.Accessory)
	}

	bridge := accessory.NewBridge(accessory.Info{
		Name:         "UDAP",
		ID:           8327,
		SerialNumber: "000-00-001",
		Manufacturer: "Braden Nicholson",
		Model:        "UDAP",
	})

	// configure the ip transport
	config := hc.Config{
		Pin:         "00100300",
		Port:        "12346",
		StoragePath: "./db",
	}

	t, err := hc.NewIPTransport(config, bridge.Accessory, accessories...)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})
	t.Start()
	return h.Config, nil
}

func (h *Homekit) Update() error {
	return nil
}

func (h *Homekit) Run() error {

	return nil
}

func init() {
	Module = Homekit{
		Config: plugin.Config{
			Name:        "homekit",
			Type:        "module",
			Description: "This module connects to homekit",
			Version:     "1.2.0",
			Author:      "Braden Nicholson",
		},
	}
}
