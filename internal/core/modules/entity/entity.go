// Copyright (c) 2022 Braden Nicholson

package entity

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

func New(db *gorm.DB) domain.EntityService {
	repo := NewRepository(db)
	service := NewService(repo)
	return service
}
