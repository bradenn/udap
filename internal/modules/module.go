// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/operators"
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/port/routes"
	"udap/internal/port/runtimes"
	"udap/internal/srv"
)

func NewModule(sys srv.System) {
	// Initialize service
	service := services.NewModuleService(
		repository.NewModuleRepository(sys.DB()),
		operators.NewModuleOperator(sys.Ctrl()))
	// Start Runtime
	sys.WithWatch(service)
	sys.Ctrl().Modules = service
	sys.WhenLoaded(func() {
		runtimes.NewModuleRuntime(service)
	})

	sys.WithRoute(routes.NewModuleRouter(service))
	// Enroll routes

}
