// Copyright (c) 2022 Braden Nicholson

package device

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

func New(db *gorm.DB) domain.DeviceService {
	repo := NewRepository(db)
	service := NewService(repo)
	return service
}
