// Copyright (c) 2022 Braden Nicholson

package attribute

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

type attributeRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.AttributeRepository {
	return &attributeRepo{
		db: db,
	}
}

func (u *attributeRepo) Register(attribute *domain.Attribute) error {
	err := u.db.Model(&domain.Attribute{}).Where(attribute).FirstOrCreate(attribute).Error
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

func (u *attributeRepo) FindAll() (*[]domain.Attribute, error) {
	var target []domain.Attribute
	if err := u.db.Find(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (u *attributeRepo) FindById(id string) (*domain.Attribute, error) {
	var target domain.Attribute
	if err := u.db.Where("id = ?", id).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (u *attributeRepo) Create(attribute *domain.Attribute) error {
	if err := u.db.Create(attribute).Error; err != nil {
		return err
	}
	return nil
}

func (u *attributeRepo) FindOrCreate(attribute *domain.Attribute) error {
	if err := u.db.FirstOrCreate(attribute).Error; err != nil {
		return err
	}
	return nil
}

func (u *attributeRepo) Update(attribute *domain.Attribute) error {
	if err := u.db.Save(attribute).Error; err != nil {
		return err
	}
	return nil
}

func (u *attributeRepo) Delete(attribute *domain.Attribute) error {
	if err := u.db.Delete(attribute).Error; err != nil {
		return err
	}
	return nil
}
