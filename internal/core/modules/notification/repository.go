// Copyright (c) 2022 Braden Nicholson

package notification

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
)

type notificationRepo struct {
	generic.Store[domain.Notification]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.NotificationRepository {
	return &notificationRepo{
		db:    db,
		Store: generic.NewStore[domain.Notification](db),
	}
}
