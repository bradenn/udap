// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/operators"
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/port/routes"
	"udap/internal/srv"
)

func NewTrigger(sys srv.System) {
	// Initialize service
	service := services.NewTriggerService(
		repository.NewTriggerRepository(sys.DB()),
		operators.NewTriggerOperator(sys.Ctrl()))
	// Enroll routes
	sys.Ctrl().Triggers = service
	sys.WithWatch(service)
	sys.WithRoute(routes.NewTriggerRouter(service))
}
