package main

import (
	"gorm.io/gorm"
	"log"
	"time"
	"udap/config"
	"udap/server"
)

func main() {
	// Load ENV variables from .env, or context
	config.Init()
	// Establish server structure, connect to database
	srv, err := server.New()
	if err != nil {
		log.Fatalln(err)
	}
	// Migrate data structures to database
	srv.Migrate(&Module{})
	srv.Migrate(&Instance{})
	srv.Migrate(&Endpoint{})
	// Configure UDAP runtime agent
	runtime := Runtime{
		endpoints: map[string]*Session{},
		server:    srv,
	}
	// Look for modules
	DiscoverModules(srv.Database())
	srv.RoutePublic("/endpoints", RouteEndpoint)
	// Register websocket handler
	srv.Router().HandleFunc("/ws", runtime.EndpointHandler)
	// Begin routes requiring jwt authentication
	// srv.RouteSecure("/instances", RouteInstances)
	// srv.RouteSecure("/modules", RouteModules)
	// Run the websocket handler server indefinitely
	go func(db *gorm.DB) {
		for {
			runtime.Update()
			time.Sleep(time.Second * 5)
		}
	}(srv.Database())
	// Run http server indefinitely
	err = srv.Run()
	if err != nil {
		log.Fatalln(err)
	}
	// If the http server exits, so too, will the websocket server and runtime agent.
}
