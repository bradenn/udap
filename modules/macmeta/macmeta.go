// Copyright (c) 2021 Braden Nicholson

package main

import (
	"os/exec"
	"strings"
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
)

var Module MacMeta

type MacMeta struct {
	plugin.Module
	localDisplay bool
	terminalId   string
	request      chan bool
	done         chan bool
}

func init() {
	config := plugin.Config{
		Name:        "macmeta",
		Type:        "module",
		Description: "MacOS meta interface for udap",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
	Module.request = make(chan bool)
	Module.done = make(chan bool)
	Module.Config = config
}

func (v *MacMeta) requestState(state bool) {
	select {
	case v.request <- state:
		return
	case <-time.After(time.Millisecond * 200):
		v.ErrF("failed to update terminal state (200ms timeout)")
		return
	}
}

func (v *MacMeta) listen() {
	for {
		select {
		case req := <-v.request:
			v.localDisplay = req
			if v.terminalId != "" {
				state := "false"
				if req {
					state = "true"
				}
				err := v.Attributes.Set(v.terminalId, "on", state)
				if err != nil {
					v.ErrF("failed to set terminal attribute: %s", err.Error())
					return
				}
			}
		case <-v.done:
			return
		}
	}
}

func (v *MacMeta) createDisplaySwitch() error {

	newSwitch := &domain.Entity{
		Name:   "terminal",
		Type:   "switch",
		Module: "macmeta",
	}
	err := v.Entities.Register(newSwitch)
	if err != nil {
		return err
	}

	v.terminalId = newSwitch.Id
	on := &domain.Attribute{
		Key:     "on",
		Value:   "true",
		Request: "true",
		Type:    "toggle",
		Channel: make(chan domain.Attribute),
		Order:   0,
		Entity:  newSwitch.Id,
	}

	go func() {
		for attribute := range on.Channel {
			if attribute.Request == "true" {
				err = v.displayOn()
				if err != nil {
					continue
				}
			} else {
				err = v.displayOff()
				if err != nil {
					continue
				}
			}
		}
	}()

	err = v.Attributes.Register(on)
	if err != nil {
		return err
	}
	return nil
}

func (v *MacMeta) displayOn() error {
	cmd := exec.Command("caffeinate", "-u", "-t", "1")
	err := cmd.Run()
	if err != nil {
		return err
	}
	v.requestState(true)
	return nil
}

func (v *MacMeta) pollDisplay() error {
	cmd := exec.Command("pmset", "-g")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	v.requestState(strings.Contains(string(output), "powerd"))
	return nil
}

func (v *MacMeta) displayOff() error {
	cmd := exec.Command("pmset", "displaysleepnow")
	err := cmd.Run()
	if err != nil {
		return err
	}
	v.requestState(false)
	return nil
}

func (v *MacMeta) Setup() (plugin.Config, error) {
	err := v.UpdateInterval(2000)
	if err != nil {
		return plugin.Config{}, err
	}
	go v.listen()
	return v.Config, nil
}

func (v *MacMeta) Update() error {
	if v.Ready() {
		err := v.pollDisplay()
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *MacMeta) Dispose() error {
	select {
	case v.done <- true:
	default:
	}
	return nil
}

func (v *MacMeta) Run() error {
	err := v.displayOn()
	if err != nil {
		return err
	}
	err = v.createDisplaySwitch()
	if err != nil {
		return err
	}
	return nil
}
