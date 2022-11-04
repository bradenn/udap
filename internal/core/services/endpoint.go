// Copyright (c) 2022 Braden Nicholson

package services

import (
	"github.com/gorilla/websocket"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

func NewEndpointService(repository ports.EndpointRepository, operator ports.EndpointOperator) ports.EndpointService {
	return &endpointService{repository: repository, operator: operator}
}

type endpointService struct {
	repository ports.EndpointRepository
	operator   ports.EndpointOperator
	generic.Watchable[domain.Endpoint]
}

func (u *endpointService) SendAll(target string, operation string, payload any) error {
	err := u.operator.SendAll(target, operation, payload)
	if err != nil {
		return err
	}
	return nil
}

func (u *endpointService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, endpoint := range *all {
		err = u.Emit(endpoint)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *endpointService) CloseAll() error {
	return u.operator.CloseAll()
}

func (u *endpointService) Send(id string, operation string, payload any) error {
	err := u.operator.Send(id, operation, payload)
	if err != nil {
		return err
	}
	return nil
}

func (u *endpointService) FindByKey(key string) (*domain.Endpoint, error) {
	return u.repository.FindByKey(key)
}

func (u *endpointService) Enroll(id string, conn *websocket.Conn) error {
	endpoint, err := u.FindById(id)
	if err != nil {
		return err
	}
	err = u.operator.Enroll(endpoint, conn)
	if err != nil {
		return err
	}
	endpoint, err = u.FindById(id)
	if err != nil {
		return err
	}
	endpoint.Connected = true
	err = u.Emit(*endpoint)
	if err != nil {
		return err
	}
	return nil
}

func (u *endpointService) Unenroll(id string) error {
	endpoint, err := u.FindById(id)
	if err != nil {
		return err
	}

	err = u.operator.Unenroll(endpoint)
	if err != nil {
		return err
	}
	endpoint.Connected = false
	err = u.Emit(*endpoint)
	if err != nil {
		return err
	}
	return nil
}

// Repository Mapping

func (u *endpointService) FindAll() (*[]domain.Endpoint, error) {
	return u.repository.FindAll()
}

func (u *endpointService) FindById(id string) (*domain.Endpoint, error) {
	return u.repository.FindById(id)
}

func (u *endpointService) Create(endpoint *domain.Endpoint) error {
	return u.repository.Create(endpoint)
}

func (u *endpointService) FindOrCreate(endpoint *domain.Endpoint) error {
	return u.repository.FindOrCreate(endpoint)
}

func (u *endpointService) Update(endpoint *domain.Endpoint) error {
	return u.repository.Update(endpoint)
}

func (u *endpointService) Delete(id string) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	return u.repository.Delete(byId)
}
