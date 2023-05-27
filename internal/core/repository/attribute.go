// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

type attributeRepo struct {
	generic.Store[domain.Attribute]
	db *gorm.DB
}

func (u *attributeRepo) FindRecentLogs() (*[]domain.AttributeLog, error) {
	var logs []domain.AttributeLog
	err := u.db.Model(&domain.AttributeLog{}).Limit(100).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return &logs, nil
}

func NewAttributeRepository(db *gorm.DB) ports.AttributeRepository {
	return &attributeRepo{
		db:    db,
		Store: generic.NewStore[domain.Attribute](db),
	}
}

func (u *attributeRepo) Log(attribute *domain.Attribute) (*domain.AttributeLog, error) {
	if attribute.Type == "media" {
		return nil, nil
	}
	log := attribute.ToLog()
	err := u.db.Model(&log).Create(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (u *attributeRepo) Register(attribute *domain.Attribute) error {
	err := u.db.Model(&domain.Attribute{}).Where("entity = ? AND key = ?", attribute.Entity, attribute.Key).FirstOrCreate(attribute).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *attributeRepo) FindByComposite(entity string, key string) (*domain.Attribute, error) {
	var target domain.Attribute
	if err := u.db.Model(&domain.Attribute{}).Where("entity = ? AND key = ?", entity,
		key).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (u *attributeRepo) FindAllByEntity(entity string) (*[]domain.Attribute, error) {
	var target []domain.Attribute
	if err := u.db.Where("entity = ?", entity).Find(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}