// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/operators"
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/port/routes"
	"udap/internal/srv"
)

func NewMacro(sys srv.System) {
	// Initialize service
	service := services.NewMacroService(
		repository.NewMacroRepository(sys.DB()),
		operators.NewMacroOperator(sys.Ctrl()))
	sys.Ctrl().Macros = service
	// Enroll routes
	sys.WithWatch(service)
	sys.WithRoute(routes.NewMacroRouter(service))
}
