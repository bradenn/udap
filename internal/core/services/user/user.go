// Copyright (c) 2022 Braden Nicholson

package user

import (
	"gorm.io/gorm"
	"udap/internal/core/ports"
	"udap/internal/core/repository"
)

func New(db *gorm.DB) ports.UserService {
	repo := repository.NewUserRepository(db)
	service := NewService(repo)
	return service
}
