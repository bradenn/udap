// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/port/routes"
	"udap/internal/srv"
)

func NewEntity(sys srv.System) {
	// Initialize service

	service := services.NewEntityService(
		repository.NewEntityRepository(sys.DB()))
	// Enroll routes

	sys.Ctrl().Entities = service

	sys.WithWatch(service)

	sys.WithRoute(routes.NewEntityRouter(service))
}
