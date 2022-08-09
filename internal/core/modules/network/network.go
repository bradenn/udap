// Copyright (c) 2022 Braden Nicholson

package network

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/repository"
)

func New(db *gorm.DB) domain.NetworkService {
	repo := repository.NewNetworkRepository(db)
	service := NewService(repo)
	return service
}
