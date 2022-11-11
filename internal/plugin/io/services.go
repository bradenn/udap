// Copyright (c) 2022 Braden Nicholson

package io

import (
	"udap/internal/core/ports"
)

type Services struct {
	Attributes ports.AttributeService
	Entities   ports.EntityService
	Triggers   ports.TriggerService
	Devices    ports.DeviceService
	Networks   ports.NetworkService
	Logs       ports.LogService
	Modules    ports.ModuleService
	Zones      ports.ZoneService
}

type ServicesIFace[T Services] interface {
	Services
}
