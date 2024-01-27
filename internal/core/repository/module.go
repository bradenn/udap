// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"gorm.io/gorm"
	"sync"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

type moduleRepo struct {
	generic.Store[domain.Module]
	cache map[string]*domain.Module
	mutex sync.Mutex
	db    *gorm.DB
}

func NewModuleRepository(db *gorm.DB) ports.ModuleRepository {
	return &moduleRepo{
		db:    db,
		Store: generic.NewStore[domain.Module](db),
		mutex: sync.Mutex{},
		cache: map[string]*domain.Module{},
	}
}

func (m *moduleRepo) FindByName(name string) (*domain.Module, error) {
	var target domain.Module
	if err := m.db.Where("name = ?", name).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (m *moduleRepo) FindByUUID(uuid string) (*domain.Module, error) {
	var target domain.Module
	if err := m.db.Where("uuid = ?", uuid).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}
