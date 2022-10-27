// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/port/routes"
	"udap/internal/srv"
)

func NewUser(sys srv.System) {
	// Initialize service
	service := services.NewUserService(
		repository.NewUserRepository(sys.DB()))
	// Enroll routes
	sys.Ctrl().Users = service
	sys.WithWatch(service)
	sys.WithRoute(routes.NewUserRouter(service))
}
