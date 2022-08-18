// Copyright (c) 2022 Braden Nicholson

package zone

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/repository"
)

func New(db *gorm.DB) domain.ZoneService {
	repo := repository.NewZoneRepository(db)
	service := NewService(repo)
	return service
}
