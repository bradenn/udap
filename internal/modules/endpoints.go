// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/operators"
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/port/routes"
	"udap/internal/srv"
)

func NewEndpoint(sys srv.System) {
	// Initialize service
	service := services.NewEndpointService(
		repository.NewEndpointRepository(sys.DB()),
		operators.NewEndpointOperator(sys.Ctrl()))
	// Enroll routes
	sys.Ctrl().Endpoints = service
	sys.WithWatch(service)
	sys.WithRoute(routes.NewEndpointRouter(service))
}
