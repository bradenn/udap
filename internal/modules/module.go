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
	sys.Ctrl().Modules = service
	err := service.Discover()
	if err != nil {
		log.Err(err)
		return
	}
	err = service.BuildAll()
	if err != nil {
		log.Err(err)
		return
	}
	err = service.LoadAll()
	if err != nil {
		return
	}

	// Start Runtime
	sys.WithWatch(service)
	sys.WithRoute(routes.NewModuleRouter(service))
	// Enroll routes
	sys.WhenLoaded(func() {
		err = service.RunAll()
		if err != nil {
			return
		}
	})

}
