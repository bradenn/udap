// Copyright (c) 2022 Braden Nicholson

package ports

import (
	"github.com/gorilla/websocket"
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type EndpointRepository interface {
	common.Persist[domain.Endpoint]
	FindByKey(key string) (*domain.Endpoint, error)
}

type EndpointOperator interface {
	Enroll(*domain.Endpoint, *websocket.Conn) error
	Unenroll(*domain.Endpoint) error
	Send(id string, operation string, payload any) error
	SendAll(id string, operation string, payload any) error
	CloseAll() error
}

type EndpointService interface {
	FindAll() (*[]domain.Endpoint, error)
	FindById(id string) (*domain.Endpoint, error)
	FindByKey(key string) (*domain.Endpoint, error)
	Create(*domain.Endpoint) error
	CloseAll() error
	domain.Observable

	Enroll(id string, conn *websocket.Conn) error

	SendAll(target string, operation string, payload any) error
	Send(id string, operation string, payload any) error
	Unenroll(key string) error

	FindOrCreate(*domain.Endpoint) error
	Update(*domain.Endpoint) error
	Delete(*domain.Endpoint) error
}
