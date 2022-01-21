// Copyright (c) 2021 Braden Nicholson

package main

import (
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"os"
	"strconv"
	"time"
	"udap/internal/models"
	"udap/pkg/plugin"
)

var Module Homekit

type Homekit struct {
	plugin.Module
}

func init() {
	config := plugin.Config{
		Name:        "homekit",
		Type:        "module",
		Description: "Homekit integration",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
}

func (h *Homekit) Setup() (plugin.Config, error) {

	return h.Config, nil
}

func (h *Homekit) Update() error {
	return nil
}

func (h *Homekit) Run() error {

	time.Sleep(time.Second * 5)

	bridge := accessory.NewBridge(accessory.Info{
		Name:             "udap",
		ID:               0120,
		SerialNumber:     "000-02-001",
		Manufacturer:     "Braden Nicholson",
		FirmwareRevision: os.Getenv("version"),
		Model:            "udap-homekit",
	})

	// configure the ip transport
	config := hc.Config{
		Pin:         "12344321",
		Port:        "34875",
		StoragePath: "./local/homekit",
	}
	var accessories []*accessory.Accessory

	compile, err := h.Entities.Compile()
	if err != nil {
		return err
	}

	for i, entity := range compile {
		switch entity.Type {
		case "spectrum":
			info := accessory.Info{
				Name:             entity.Name,
				ID:               uint64(i + 2),
				SerialNumber:     entity.Id,
				Manufacturer:     fmt.Sprintf("udap:%s", entity.Module),
				Model:            fmt.Sprintf("udap:%s", entity.Type),
				FirmwareRevision: "0.0.1",
			}
			device := accessory.NewColoredLightbulb(info)

			device.Lightbulb.Hue.OnValueRemoteGet(func() float64 {
				attr := models.Attribute{}
				attr.Entity = entity.Id
				attr.Key = "hue"
				a := h.Attributes.Find(entity.Path())
				parseInt, err := strconv.ParseInt(a.Value, 10, 64)
				if err != nil {
					return 0
				}
				return float64(parseInt) / 360.0
			})

			device.Lightbulb.Hue.OnValueRemoteUpdate(func(value float64) {
				attr := models.Attribute{}
				attr.Entity = entity.Id
				attr.Key = "hue"
				err := h.Attributes.Update(entity.Id, "hue", fmt.Sprintf("%d", int(value*360.0)))
				if err != nil {

				}
				return
			})

			device.Lightbulb.On.OnValueRemoteGet(func() bool {
				attr := models.Attribute{}
				attr.Entity = entity.Id
				attr.Key = "on"
				a := h.Attributes.Find(attr.Path())
				value, err := strconv.ParseBool(a.Value)
				if err != nil {
					return false
				}
				return value
			})

			device.Lightbulb.On.OnValueRemoteUpdate(func(value bool) {
				attr := models.Attribute{}
				attr.Entity = entity.Id
				attr.Key = "on"
				str := "false"
				if value {
					str = "true"
				}
				err := h.Attributes.Update(entity.Id, "on", str)
				if err != nil {

				}
				return
			})

			device.Lightbulb.Brightness.OnValueRemoteGet(func() int {
				attr := models.Attribute{}
				attr.Entity = entity.Id
				attr.Key = "on"
				a := h.Attributes.Find(attr.Path())
				parseInt, err := strconv.ParseInt(a.Value, 10, 64)
				if err != nil {
					return 0
				}
				return int(parseInt)
			})

			device.Lightbulb.Brightness.OnValueRemoteUpdate(func(value int) {
				attr := models.Attribute{}
				attr.Entity = entity.Id
				attr.Key = "dim"
				err := h.Attributes.Update(entity.Id, "dim", fmt.Sprintf("%d", value))
				if err != nil {

				}
				return
			})

			accessories = append(accessories, device.Accessory)
			break
		default:
			break
		}
	}

	t, err := hc.NewIPTransport(config, bridge.Accessory, accessories...)
	if err != nil {
		return err
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Handle(func() {

	})

	t.Start()
	return nil
}
