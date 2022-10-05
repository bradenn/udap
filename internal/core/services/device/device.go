// Copyright (c) 2022 Braden Nicholson

package device

import (
	"gorm.io/gorm"
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func New(db *gorm.DB) ports.DeviceService {
	repo := repository.NewDeviceRepository(db)
	service := NewService(repo)
	return service
}
