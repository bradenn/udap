// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

type attributeRepo struct {
	generic.Store[domain.Attribute]
	cache map[string]domain.Attribute
	mutex sync.Mutex
	db    *gorm.DB
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
		mutex: sync.Mutex{},
		cache: make(map[string]domain.Attribute),
	}
}

func (u *attributeRepo) Log(attribute *domain.Attribute) (*domain.AttributeLog, error) {
	if attribute.Type == "media" {
		return nil, nil
	}
	//log := attribute.ToLog()
	//err := u.db.Model(&log).Create(&log).Error
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}

func (u *attributeRepo) StateUpdate(attribute *domain.Attribute) error {
	composite := fmt.Sprintf("%s.%s", attribute.Entity, attribute.Key)

	u.mutex.Lock()
	val, ok := u.cache[composite]
	u.mutex.Unlock()

	if !ok {
		u.cache[composite] = *attribute
	}

	if val.Value != attribute.Value {
		err := u.Store.Update(attribute)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *attributeRepo) CachedValue(entity string, key string) (*domain.Attribute, error) {
	composite := fmt.Sprintf("%s.%s", entity, key)

	u.mutex.Lock()
	val, ok := u.cache[composite]
	u.mutex.Unlock()

	if ok {
		return &val, nil
	}

	var target domain.Attribute
	if err := u.db.Model(&domain.Attribute{}).Where("entity = ? AND key = ?", entity,
		key).First(&target).Error; err != nil {
		return nil, err
	}

	u.cache[composite] = target

	return &target, nil
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
