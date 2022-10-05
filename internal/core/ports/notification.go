// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type NotificationRepository interface {
	common.Persist[domain.Notification]
}

type NotificationOperator interface {
	Send(notification domain.Notification) error
}

type NotificationService interface {
	domain.Observable
	FindAll() (*[]domain.Notification, error)
	FindById(id string) (*domain.Notification, error)
	Create(*domain.Notification) error
	FindOrCreate(*domain.Notification) error
	Register(*domain.Notification) error
	Update(*domain.Notification) error
	Delete(*domain.Notification) error
}
