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

func (s *subRoutineRepo) FindByTriggerId(id string) (res []*domain.SubRoutine, err error) {
	err = s.db.Model(&domain.SubRoutine{}).Where("triggerId = ?", id).Find(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
