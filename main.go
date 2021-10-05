package main

import (
	"fmt"
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

	mod := Module{}
	err = srv.Database().Model(&Module{}).Where("id = ?", "c9cdff54-025a-47dc-abc5-1dc627f43894").First(&mod).Error
	if err != nil {
		return
	}

	instance, err := mod.CreateInstance(srv.Database(), "Braden's Spotify")
	if err != nil {
		return
	}
	fmt.Println(instance)

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
