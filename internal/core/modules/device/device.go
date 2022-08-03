// Copyright (c) 2022 Braden Nicholson

package device

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/repository"
)

func New(db *gorm.DB) domain.DeviceService {
	repo := repository.NewDeviceRepository(db)
	service := NewService(repo)
	return service
}
