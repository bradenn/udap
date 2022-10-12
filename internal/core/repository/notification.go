// Copyright (c) 2022 Braden Nicholson

package repository

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

type notificationRepo struct {
	generic.Store[domain.Notification]
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) ports.NotificationRepository {
	return &notificationRepo{
		db:    db,
		Store: generic.NewStore[domain.Notification](db),
	}
}
