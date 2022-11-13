// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/srv"
)

func NewNetwork(sys srv.System) {
	// Initialize service
	service := services.NewNetworkService(
		repository.NewNetworkRepository(sys.DB()))
	// Enroll routes
	sys.WithWatch(service)
	// sys.WithRoute(routes.NewNetworkRouter(service))
	sys.Ctrl().Networks = service
}
