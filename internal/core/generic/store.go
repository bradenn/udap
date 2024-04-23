// Copyright (c) 2022 Braden Nicholson

package generic

import (
	"gorm.io/gorm"
	"sync"
	"udap/internal/core/domain"
)

type PersistentType interface {
	domain.User | domain.Module | domain.Entity | domain.Device | domain.Attribute | domain.Endpoint | domain.
		Network | domain.Zone | domain.Notification | domain.Macro | domain.Trigger | domain.SubRoutine | domain.AttributeLog | Mock
}

type Store[T any] struct {
	cache     map[string]T
	cacheLock sync.RWMutex
	db        *gorm.DB
}

func NewStore[T any](db *gorm.DB) Store[T] {
	return Store[T]{
		db:        db,
		cache:     map[string]T{},
		cacheLock: sync.RWMutex{},
	}
}

// FindAll returns all records of the type T
func (c *Store[T]) FindAll() (*[]T, error) {
	var target []T
	if err := c.db.Model(&target).Find(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

// FindById returns the first record with a UUID matching the provided string
func (c *Store[T]) FindById(id string) (*T, error) {
	var target T
	if err := c.db.Model(&target).Where("id = ?", id).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

// Create creates a record of the type T
func (c *Store[T]) Create(t *T) error {
	if err := c.db.Model(&t).Create(t).Error; err != nil {
		return err
	}
	return nil
}

// FindOrCreate will emplace any record into its appropriate table
func (c *Store[T]) FindOrCreate(t *T) error {
	if err := c.db.Model(&t).FirstOrCreate(t).Error; err != nil {
		return err
	}
	return nil
}

// Update saves any changes made to the provided record of type T
func (c *Store[T]) Update(t *T) error {
	if err := c.db.Model(&t).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes the record from the database
func (c *Store[T]) Delete(t *T) error {
	if err := c.db.Model(&t).Delete(t).Error; err != nil {
		return err
	}
	return nil
}
