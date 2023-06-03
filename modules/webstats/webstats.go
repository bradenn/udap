// Copyright (c) 2021 Braden Nicholson

package main

import (
	"time"
	"udap/internal/plugin"
)

var Module WebStats

type WebStats struct {
	plugin.Module
	eId string
}

func init() {
	config := plugin.Config{
		Name:        "webstats",
		Type:        "module",
		Description: "Web related statistics",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
}

func (w *WebStats) Setup() (plugin.Config, error) {
	err := w.UpdateInterval(2000)
	if err != nil {
		return plugin.Config{}, err
	}
	return w.Config, nil
}

func (w *WebStats) pull() error {
	time.Sleep(250 * time.Millisecond)
	return nil
}

func (w *WebStats) Update() error {
	if w.Ready() {
		return w.pull()
	}
	return nil
}

func (w *WebStats) Run() error {
	return nil
}