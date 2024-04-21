// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"time"
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type AttributeRepository interface {
	common.Persist[domain.Attribute]
	FindAllByEntity(entity string) (*[]domain.Attribute, error)
	FindByComposite(entity string, key string) (*domain.Attribute, error)
	Log(attribute *domain.Attribute) error
	Register(*domain.Attribute) error
	StateUpdate(*domain.Attribute) error
	FindRecentLogs() (*[]domain.AttributeLog, error)
	FindRecent() (*[]domain.Attribute, error)
	Summary(key string, start int64, stop int64, window int, mode string) (map[int64]float64, error)
}

type AttributeOperator interface {
	Register(attribute *domain.Attribute) error
	Request(*domain.Attribute, string) error
	Set(*domain.Attribute, string) error
	Update(*domain.Attribute, string, time.Time) error
}

type AttributeService interface {
	domain.Observable
	FindAll() (*[]domain.Attribute, error)
	FindByComposite(entity string, key string) (*domain.Attribute, error)
	FindAllByEntity(entity string) (*[]domain.Attribute, error)
	FindById(id string) (*domain.Attribute, error)
	Create(*domain.Attribute) error
	Register(*domain.Attribute) error
	Summary(key string, start int64, stop int64, window int, mode string) (map[int64]float64, error)
	Request(entity string, key string, value string) error
	Set(entity string, key string, value string) error
	Update(entity string, key string, value string, stamp time.Time) error
	Delete(*domain.Attribute) error
}
