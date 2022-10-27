// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/port/routes"
	"udap/internal/srv"
)

func NewDevice(sys srv.System) {
	// Initialize service
	service := services.NewDeviceService(
		repository.NewDeviceRepository(sys.DB()))
	// Enroll routes
	sys.WithWatch(service)
	sys.WithRoute(routes.NewDeviceRouter(service))
	sys.Ctrl().Devices = service
}
