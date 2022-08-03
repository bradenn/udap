// Copyright (c) 2022 Braden Nicholson

package notification

import (
	"gorm.io/gorm"
	"udap/internal/core/domain"
	"udap/internal/core/repository"
)

func New(db *gorm.DB) domain.NotificationService {
	repo := repository.NewNotificationRepository(db)
	service := NewService(repo)
	return service
}
