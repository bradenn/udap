// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/operators"
	"udap/internal/core/repository"
	"udap/internal/core/services"
	"udap/internal/port/routes"
	"udap/internal/srv"
)

func NewSubroutine(sys srv.System) {
	// Initialize service
	service := services.NewSubRoutineService(
		repository.NewSubRoutineRepository(sys.DB()),
		operators.NewSubRoutineOperator(sys.Ctrl()))
	// Enroll routes
	sys.Ctrl().SubRoutines = service
	sys.WithWatch(service)
	sys.WithRoute(routes.NewSubroutineRouter(service))
}
