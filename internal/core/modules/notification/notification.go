// Copyright (c) 2022 Braden Nicholson

package notification

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
)

func New(db *gorm.DB) domain.NotificationService {
	repo := NewRepository(db)
	service := NewService(repo)
	return service
}
