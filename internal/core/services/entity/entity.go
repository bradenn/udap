// Copyright (c) 2022 Braden Nicholson

package entity

import (
	"gorm.io/gorm"
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func New(db *gorm.DB) ports.EntityService {
	repo := repository.NewEntityRepository(db)
	service := NewService(repo)
	return service
}
