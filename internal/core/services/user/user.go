// Copyright (c) 2022 Braden Nicholson

package user

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/repository"
)

func New(db *gorm.DB) domain.UserService {
	repo := repository.NewUserRepository(db)
	service := NewService(repo)
	return service
}
