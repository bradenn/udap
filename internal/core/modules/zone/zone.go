// Copyright (c) 2022 Braden Nicholson

package zone

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

func New(db *gorm.DB) domain.ZoneService {
	repo := NewRepository(db)
	service := NewService(repo)
	return service
}
