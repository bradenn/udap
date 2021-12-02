// Copyright (c) 2021 Braden Nicholson

package main

import (
	"fmt"
	"udap/pkg/plugin"
)

var Plugin Sample

var P plugin.Module

func init() {
	config := plugin.Config{
		Name:        "Sample",
		Type:        "module",
		Description: "sample",
		Version:     "0.1.0",
		Author:      "Braden Nicholson",
	}
	plugin.NewModule(&P, config)

}

type Sample struct {
	requests chan plugin.Request
	resolver chan plugin.Event
	metadata plugin.Metadata
}

func (s *Sample) Startup() (plugin.Metadata, error) {
	Plugin = Sample{
		metadata: plugin.Metadata{
			Name:        "Sample",
			Type:        "module",
			Description: "This is a sample plugin",
			Version:     "1.1.0",
			Author:      "BN",
		},
		requests: make(chan plugin.Request),
	}
	return s.metadata, nil
}

func (s *Sample) Metadata() plugin.Metadata {
	return s.metadata
}

func (s *Sample) Listen() {
	for request := range s.requests {
		fmt.Println(request)
	}
	close(s.requests)
}

func (s *Sample) Connect(events chan plugin.Event) chan plugin.Request {
	s.resolver = events
	return s.requests
}

func (s *Sample) resolve(event plugin.Event) {
	s.resolver <- event
}

func (s *Sample) Cleanup() {
	close(s.requests)
}
