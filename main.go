package main

import (
	"time"
	"udap/config"
	"udap/logger"
	"udap/server"
)

func main() {
	// Load ENV variables from .env, or context
	config.Init()
	// Establish server structure, connect to database
	srv, err := server.New()
	if err != nil {
		logger.Err(err)
		return
	}

	// Migrate data structures to database
	srv.Migrate(&Module{})
	srv.Migrate(&Instance{})
	srv.Migrate(&Endpoint{})
	// Configure UDAP runtime agent
	runtime := NewRuntime(srv)
	runtime.Begin(time.Second * 5)
	// Route Endpoint authentication
	srv.RoutePublic("/endpoints", RouteEndpoint)
	// Run http server indefinitely
	err = srv.Run()
	if err != nil {
		logger.Err(err)
	}
	// If the http server exits, so too, will the websocket server and runtime agent.
}
