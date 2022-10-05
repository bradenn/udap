// Copyright (c) 2022 Braden Nicholson

package ports

import "udap/internal/core/domain"

type LogRepository interface {
	FindAll() ([]domain.Log, error)
	Create(*domain.Log) error
}

type LogService interface {
	domain.Observable
	Create(*domain.Log) error
}
