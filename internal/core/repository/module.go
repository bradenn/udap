// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
)

type moduleRepo struct {
	generic.Store[domain.Module]
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) domain.ModuleRepository {
	return &moduleRepo{
		db:    db,
		Store: generic.NewStore[domain.Module](db),
	}
}

func (m moduleRepo) FindByName(name string) (*domain.Module, error) {
	var target domain.Module
	if err := m.db.Where("name = ?", name).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}
