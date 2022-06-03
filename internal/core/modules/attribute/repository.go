// Copyright (c) 2022 Braden Nicholson

package attribute

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
)

type attributeRepo struct {
	generic.Store[domain.Attribute]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.AttributeRepository {
	return &attributeRepo{
		db:    db,
		Store: generic.NewStore[domain.Attribute](db),
	}
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
