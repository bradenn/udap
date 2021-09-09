package main

import (
	"log"
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
	// Begin routes requiring jwt authentication
	srv.RouteSecure("/endpoints", &Endpoint{})
	srv.RouteSecure("/instances", &Instance{})
	// Run the server indefinitely
	err = srv.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
