// Copyright (c) 2022 Braden Nicholson

package module

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

type moduleRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.ModuleRepository {
	return &moduleRepo{
		db: db,
	}
}

func (m moduleRepo) FindByName(name string) (*domain.Module, error) {
	var target domain.Module
	if err := m.db.Where("name = ?", name).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (m moduleRepo) FindAll() (*[]domain.Module, error) {
	var targets []domain.Module
	if err := m.db.Find(&targets).Error; err != nil {
		return nil, err
	}
	return &targets, nil
}

func (m moduleRepo) FindById(id string) (*domain.Module, error) {
	var target *domain.Module
	if err := m.db.Where("id = ?", id).First(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (m moduleRepo) Create(module *domain.Module) error {
	if err := m.db.Create(module).Error; err != nil {
		return err
	}
	return nil
}

func (m moduleRepo) FindOrCreate(module *domain.Module) error {
	if err := m.db.FirstOrCreate(module).Error; err != nil {
		return err
	}
	return nil
}

func (m moduleRepo) Update(module *domain.Module) error {
	if err := m.db.Save(module).Error; err != nil {
		return err
	}
	return nil
}

func (m moduleRepo) Delete(module *domain.Module) error {
	if err := m.db.Delete(module).Error; err != nil {
		return err
	}
	return nil
}
