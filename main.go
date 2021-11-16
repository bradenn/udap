package main

import (
	"udap/config"
	"udap/runtime"
	"udap/server"
	"udap/types"
)

func main() {
	// Load ENV variables from .env, or context
	config.Init()
	// Establish server structure, connect to database
	srv, err := server.New()
	if err != nil {
		config.Err(err)
		return
	}
	types.Load(srv.Database)
	// Migrate data structures to database
	srv.Migrate(&types.Module{})
	srv.Migrate(&types.Instance{})
	srv.Migrate(&types.Endpoint{})
	srv.Migrate(&types.Entity{})
	// Configure UDAP runtime agent
	runtime.New(srv)
	// Route Endpoint authentication
	srv.RoutePublic("/endpoints", types.RouteEndpoint)
	// Run http server indefinitely
	err = srv.Run()
	if err != nil {
		config.Err(err)
	}
	// If the http server exits, so too, will the websocket server and runtime agent.
}
