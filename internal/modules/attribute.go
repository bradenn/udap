// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/operators"
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/port/routes"
	"udap/internal/srv"
)

func NewAttribute(sys srv.System) {
	// Initialize service
	service := services.NewAttributeService(
		repository.NewAttributeRepository(sys.DB()),
		operators.NewAttributeOperator())
	// Enroll routes
	sys.WithWatch(service)
	sys.WithRoute(routes.NewAttributeRouter(service))
	sys.Ctrl().Attributes = service
}
