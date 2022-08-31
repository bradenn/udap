// Copyright (c) 2022 Braden Nicholson

package entity

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/repository"
)

func New(db *gorm.DB) domain.EntityService {
	repo := repository.NewEntityRepository(db)
	service := NewService(repo)
	return service
}
