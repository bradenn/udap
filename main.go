// Copyright (c) 2021 Braden Nicholson

package main

import (
	"udap/endpoint"
	"udap/module"
	"udap/runtime"
	"udap/server"
	"udap/types"
	"udap/udap"
	"udap/udap/db"
)

func main() {

	// Load ENV variables from .env, or context
	udap.Init()
	_, err := db.NewGormDB()
	if err != nil {
		return
	}
	// Establish server structure, connect to database

	srv, err := server.New()
	if err != nil {
		udap.Err(err)
		return
	}
	types.Load(db.DB)
	udap.New()
	// Migrate data structures to database
	srv.Migrate(&module.Module{})
	srv.Migrate(&module.Instance{})
	srv.Migrate(&endpoint.Endpoint{})
	srv.Migrate(&types.Entity{})
	srv.Migrate(&endpoint.Subscription{})
	srv.Migrate(&endpoint.Grant{})
	// Configure UDAP runtime agent
	runtime.New(srv)
	// Route Endpoint authentication
	srv.RoutePublic("/endpoints", endpoint.RouteEndpoint)
	// Run http server indefinitely
	err = srv.Run()
	if err != nil {
		udap.Err(err)
	}
	// If the http server exits, so too, will the websocket server and runtime agent.
}
