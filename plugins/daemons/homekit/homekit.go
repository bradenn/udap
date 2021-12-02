// Copyright (c) 2021 Braden Nicholson

package main

import (
	"fmt"
	"udap/pkg/plugin"
)

var Plugin Homekit

type Homekit struct {
	requests chan plugin.Request
	resolver chan plugin.Event
	metadata plugin.Metadata
}

func (s *Homekit) Startup() (plugin.Metadata, error) {
	Plugin = Homekit{
		metadata: plugin.Metadata{
			Name:        "Homekit",
			Type:        "daemon",
			Description: "This plugin connects to homekit",
			Version:     "1.2.0",
			Author:      "Braden Nicholson",
		},
		requests: make(chan plugin.Request),
	}
	return s.metadata, nil
}

func (s *Homekit) Metadata() plugin.Metadata {
	return s.metadata
}

func (s *Homekit) Listen() {
	for request := range s.requests {
		fmt.Println(request)
	}
	close(s.requests)
}

func (s *Homekit) Connect(events chan plugin.Event) chan plugin.Request {
	s.resolver = events
	go s.Listen()
	return s.requests
}

func (s *Homekit) resolve(event plugin.Event) {
	s.resolver <- event
}

func (s *Homekit) Cleanup() {
	close(s.requests)
}
