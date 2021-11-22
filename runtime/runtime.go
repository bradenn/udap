// Copyright (c) 2021 Braden Nicholson

package runtime

import (
	"udap/endpoint"
	"udap/module"
	"udap/server"
)

// Runtime manages the event-loop for all instances, as well as the websocket connections between UDAP and endpoints.
type Runtime struct {
	modules map[string]*module.Module
	// endpoints is a map of each active session
	endpoints map[string]*endpoint.Endpoint
	// instances contains all active instances of udap
	instances map[string]*module.Instance
	// cache holds the current data for each live instance
	cache map[string]string

	updater chan module.UpdateBuffer

	server server.Server
}

// New creates, initializes, and configures a websocket runtime
func New(server server.Server) Runtime {
	// The updater channel is called by modules, we initialize it to queue 16 at a time
	updater := make(chan module.UpdateBuffer, 16)
	runtime := Runtime{
		modules:   map[string]*module.Module{},
		endpoints: map[string]*endpoint.Endpoint{},
		instances: map[string]*module.Instance{},
		cache:     map[string]string{},
		updater:   updater,
		server:    server,
	}

	// Discover any unloaded modules
	runtime.discoverModules()
	// Route WebSocket requests here
	server.Router().HandleFunc("/ws/{token}", endpoint.HandleSockets)
	// Set async listener for the updating from modules and daemons
	go runtime.Listen()
	// return the runtime
	return runtime
}
