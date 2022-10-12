// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

type triggerRepo struct {
	generic.Store[domain.Trigger]
	db *gorm.DB
}

func NewTriggerRepository(db *gorm.DB) ports.TriggerRepository {
	return &triggerRepo{
		db:    db,
		Store: generic.NewStore[domain.Trigger](db),
	}
}
