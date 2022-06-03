// Copyright (c) 2022 Braden Nicholson

package network

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

func New(db *gorm.DB) domain.NetworkService {
	repo := NewRepository(db)
	service := NewService(repo)
	return service
}
