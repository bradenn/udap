// Copyright (c) 2022 Braden Nicholson

package notification

import (
	"fmt"
	"udap/internal/core/domain"
)

type notificationService struct {
	repository domain.NotificationRepository
	channel    chan<- domain.Mutation
}

func (u *notificationService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, notification := range *all {
		err = u.emit(&notification)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *notificationService) emit(notification *domain.Notification) error {
	if u.channel == nil {
		return nil
	}
	u.channel <- domain.Mutation{
		Status:    "update",
		Operation: "notification",
		Body:      *notification,
		Id:        notification.Id,
	}
	return nil
}

func (u *notificationService) Watch(mut chan<- domain.Mutation) error {
	if u.channel != nil {
		return fmt.Errorf("channel already set")
	}
	u.channel = mut
	return nil
}

func (u notificationService) Register(notification *domain.Notification) error {
	return nil
}

func NewService(repository domain.NotificationRepository) domain.NotificationService {
	return &notificationService{repository: repository}
}

// Repository Mapping

func (u notificationService) FindAll() (*[]domain.Notification, error) {
	return u.repository.FindAll()
}

func (u notificationService) FindById(id string) (*domain.Notification, error) {
	return u.repository.FindById(id)
}

func (u notificationService) Create(notification *domain.Notification) error {
	return u.repository.Create(notification)
}

func (u notificationService) FindOrCreate(notification *domain.Notification) error {
	return u.repository.FindOrCreate(notification)
}

func (u notificationService) Update(notification *domain.Notification) error {
	return u.repository.Update(notification)
}

func (u notificationService) Delete(notification *domain.Notification) error {
	return u.repository.Delete(notification)
}
