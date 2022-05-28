// Copyright (c) 2022 Braden Nicholson

package module

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

func New(db *gorm.DB) domain.ModuleService {
	repo := NewRepository(db)
	service := NewService(repo)
	return service
}
