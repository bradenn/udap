// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/operators"
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/log"
	"udap/internal/port/routes"
	"udap/internal/srv"
)

func NewModule(sys srv.System) {
	// Initialize service
	service := services.NewModuleService(
		repository.NewModuleRepository(sys.DB()),
		operators.NewModuleOperator(sys.Ctrl()))
	// Publish the service
	sys.Ctrl().Modules = service
	// Assign mutation channel
	sys.WithWatch(service)
	// Discover local modules
	err := service.Discover()
	if err != nil {
		log.Err(err)

	}

	// Build local modules
	err = service.BuildAll()
	if err != nil {
		log.Err(err)
	}

	// Load all modules
	err = service.LoadAll()
	if err != nil {
		log.Err(err)
	}
	// Assign routes

	sys.WithRoute(routes.NewModuleRouter(service))
	// Start all modules

	sys.WhenLoaded(func() {
		log.Event("Starting modules...")
		err = service.RunAll()
		if err != nil {
			return
		}
	})

}
