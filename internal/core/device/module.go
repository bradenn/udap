// Copyright (c) 2022 Braden Nicholson

package device

import (
	"time"
	"udap/internal/core/ports"
	"udap/internal/srv"
)

type ModuleMeta struct {
	loadedAt  time.Time
	updatedAt time.Time
}

type Module struct {
	ModuleMeta
	repository ports.DeviceRepository
	service    ports.DeviceService
}

func NewModule(rtx srv.System) {
	repository := newRepository(rtx.DB())
	// Initialize service
	service := newService(repository)
	// Enroll routes
	rtx.WithWatch(service)
	rtx.WithRoute(newRouter(service))
	rtx.Ctrl().Devices = service
}
