// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

type triggerRepo struct {
	generic.Store[domain.User]
	db *gorm.DB
}

func NewTriggerRepository(db *gorm.DB) ports.UserRepository {
	return &userRepo{
		db:    db,
		Store: generic.NewStore[domain.User](db),
	}
}
