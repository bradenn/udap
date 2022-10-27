// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/port/routes"
	"udap/internal/srv"
)

func NewZone(sys srv.System) {
	// Initialize service
	service := services.NewZoneService(
		repository.NewZoneRepository(sys.DB()))
	// Enroll routes
	sys.WithWatch(service)
	sys.WithRoute(routes.NewZoneRouter(service))
	sys.Ctrl().Zones = service
}
