// Copyright (c) 2022 Braden Nicholson

package modules

import (
	"udap/internal/core/services"
	"udap/internal/srv"
)

func NewLog(sys srv.System) {
	// Initialize service
	service := services.NewLogService()
	sys.Ctrl().Logs = service
	// Enroll routes
	sys.WithWatch(service)
}
