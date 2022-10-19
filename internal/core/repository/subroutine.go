// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

type subRoutineRepo struct {
	generic.Store[domain.SubRoutine]
	db *gorm.DB
}

func NewSubRoutineRepository(db *gorm.DB) ports.SubRoutineRepository {
	return &subRoutineRepo{
		db:    db,
		Store: generic.NewStore[domain.SubRoutine](db),
	}
}

func (s *subRoutineRepo) FindById(id string) (*domain.SubRoutine, error) {
	target := domain.SubRoutine{}
	err := s.db.Model(&domain.SubRoutine{}).Preload("Macros").Where("id = ?", id).Find(&target).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (s *subRoutineRepo) FindAll() (*[]domain.SubRoutine, error) {
	var target []domain.SubRoutine
	if err := s.db.Preload("Macros").Find(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (s *subRoutineRepo) FindByTriggerId(id string) (res []*domain.SubRoutine, err error) {
	var rs []*domain.SubRoutine
	err = s.db.Model(&domain.SubRoutine{}).Preload("Macros").Where("trigger_id = ?", id).Find(&rs).Error
	if err != nil {
		return nil, err
	}
	return rs, nil
}
