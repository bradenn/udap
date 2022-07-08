// Copyright (c) 2022 Braden Nicholson

package domain

import "udap/internal/core/domain/common"

type Notification struct {
	common.Persistent
	Title    string `json:"title"`
	Target   string `json:"string"`
	Module   string `json:"module"`
	Body     string `json:"body"`
	Priority int    `json:"priority"`
}

type NotificationRepository interface {
	common.Persist[Notification]
}

type NotificationOperator interface {
	Send(notification Notification) error
}

type NotificationService interface {
	Observable
	FindAll() (*[]Notification, error)
	FindById(id string) (*Notification, error)
	Create(*Notification) error
	FindOrCreate(*Notification) error
	Register(*Notification) error
	Update(*Notification) error
	Delete(*Notification) error
}
