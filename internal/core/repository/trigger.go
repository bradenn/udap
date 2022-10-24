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

func (t *triggerRepo) Register(trigger *domain.Trigger) error {
	err := t.db.Model(&domain.Trigger{}).Where("name = ?", trigger.Name).FirstOrCreate(trigger).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *triggerRepo) FindByName(name string) (*domain.Trigger, error) {
	trigger := domain.Trigger{}
	err := t.db.Model(domain.Trigger{}).Where("name = ?", name).First(&trigger).Error
	return &trigger, err
}
