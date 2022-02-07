// Copyright (c) 2021 Braden Nicholson

package main

import (
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
	"os"
	"time"
	"udap/internal/controller"
	"udap/internal/log"
	"udap/internal/models"
	"udap/pkg/plugin"
)

var Module Homekit

type Homekit struct {
	plugin.Module
	bridge  *accessory.Bridge
	config  hc.Config
	devices map[string]*service.Service
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
	h.devices = map[string]*service.Service{}
	h.bridge = accessory.NewBridge(accessory.Info{
		Name:             "udap",
		ID:               1,
		SerialNumber:     "000-02-001",
		Manufacturer:     "Braden Nicholson",
		FirmwareRevision: os.Getenv("version"),
		Model:            "udap",
	})

	// configure the ip transport
	h.config = hc.Config{
		Pin:         "00102003",
		Port:        "12345",
		StoragePath: "./local/homekit",
	}
	return h.Config, nil
}

func (h *Homekit) Update() error {
	return nil
}

func (h *Homekit) Run() error {

	time.Sleep(time.Second * 5)

	var accessories []*accessory.Accessory

	keys := h.Entities.Keys()

	for _, name := range keys {
		entity := *h.Entities.Find(name)
		switch entity.Type {
		case "spectrum":
			info := accessory.Info{
				Name:             entity.Name,
				ID:               uint64(entity.CreatedAt.UnixNano()),
				SerialNumber:     entity.Id,
				Manufacturer:     fmt.Sprintf("%s", entity.Module),
				Model:            fmt.Sprintf("%s", entity.Type),
				FirmwareRevision: h.Module.Version,
			}
			device := newSpectrumLight(info)
			err := device.syncAttributes(h.Attributes, entity.Id)
			if err != nil {
				return err
			}
			accessories = append(accessories, device.Accessory)
		case "switch":
			info := accessory.Info{
				Name:             entity.Name,
				ID:               uint64(entity.CreatedAt.UnixNano()),
				SerialNumber:     entity.Id,
				Manufacturer:     fmt.Sprintf("%s", entity.Module),
				Model:            fmt.Sprintf("%s", entity.Type),
				FirmwareRevision: h.Module.Version,
			}
			device := accessory.NewSwitch(info)
			syncSwitch(device.Switch, h.Attributes, entity.Id)

			accessories = append(accessories, device.Accessory)
		default:

		}

	}

	t, err := hc.NewIPTransport(h.config, h.bridge.Accessory, accessories...)
	if err != nil {
		log.Err(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()

	return nil
}

func syncSwitch(p *service.Switch, a *controller.Attributes, id string) {
	p.On.OnValueRemoteUpdate(func(b bool) {
		str := "false"
		if b {
			str = "true"
		}
		err := a.Request(id, "on", str)
		if err != nil {
			log.Err(err)
		}
	})

	a.WatchSingle(fmt.Sprintf("%s.%s", id, "on"), func(data any) error {
		attr := *data.(*models.Attribute)
		p.On.UpdateValue(attr.Request)
		return nil
	})

}

type spectrumLight struct {
	*accessory.Accessory
	spectrum *spectrum
}

func (s *spectrumLight) syncAttributes(a *controller.Attributes, id string) error {
	s.spectrum.On.OnValueRemoteUpdate(func(b bool) {
		str := "false"
		if b {
			str = "true"
		}
		err := a.Request(id, "on", str)
		if err != nil {
			log.Err(err)
		}
	})

	s.spectrum.Dim.OnValueRemoteUpdate(func(value int) {
		err := a.Request(id, "dim", fmt.Sprintf("%d", int(value/5)*5))
		if err != nil {
			log.Err(err)
		}
	})

	a.WatchSingle(fmt.Sprintf("%s.%s", id, "on"), func(data any) error {
		attr := *data.(*models.Attribute)
		s.spectrum.On.UpdateValue(attr.Request)
		return nil
	})

	a.WatchSingle(fmt.Sprintf("%s.%s", id, "dim"), func(data any) error {
		attr := *data.(*models.Attribute)
		s.spectrum.Dim.UpdateValue(attr.Request)
		return nil
	})

	return nil
}

func newSpectrumLight(info accessory.Info) *spectrumLight {
	acc := spectrumLight{}
	acc.Accessory = accessory.New(info, accessory.TypeLightbulb)
	acc.spectrum = newSpectrum()

	acc.spectrum.Dim.SetValue(100)

	acc.AddService(acc.spectrum.Service)

	return &acc
}

type spectrum struct {
	*service.Service

	On  *characteristic.On
	Dim *characteristic.Brightness
	Cct *characteristic.ColorTemperature
}

func newSpectrum() *spectrum {
	svc := spectrum{}
	svc.Service = service.New(service.TypeLightbulb)

	svc.On = characteristic.NewOn()
	svc.AddCharacteristic(svc.On.Characteristic)

	svc.Dim = characteristic.NewBrightness()
	svc.AddCharacteristic(svc.Dim.Characteristic)

	return &svc
}
