// Copyright (c) 2022 Braden Nicholson

package endpoint

import (
	"github.com/gorilla/websocket"
	"udap/internal/core/domain"
)

type endpointService struct {
	repository domain.EndpointRepository
	operator   domain.EndpointOperator
}

func (u endpointService) Send(id string, operation string, payload any) error {
	err := u.operator.Send(id, operation, payload)
	if err != nil {
		return err
	}
	return nil
}

func (u endpointService) FindByKey(key string) (*domain.Endpoint, error) {
	return u.repository.FindByKey(key)
}

func (u endpointService) Enroll(id string, conn *websocket.Conn) error {
	endpoint, err := u.FindById(id)
	if err != nil {
		return err
	}
	err = u.operator.Enroll(endpoint, conn)
	if err != nil {
		return err
	}
	return nil
}

func (u endpointService) Disconnect(key string) error {
	// TODO implement me
	panic("implement me")
}

func NewService(repository domain.EndpointRepository, operator domain.EndpointOperator) domain.EndpointService {
	return endpointService{repository: repository, operator: operator}
}

// Repository Mapping

func (u endpointService) FindAll() (*[]domain.Endpoint, error) {
	return u.repository.FindAll()
}

func (u endpointService) FindById(id string) (*domain.Endpoint, error) {
	return u.repository.FindById(id)
}

func (u endpointService) Create(endpoint *domain.Endpoint) error {
	return u.repository.Create(endpoint)
}

func (u endpointService) FindOrCreate(endpoint *domain.Endpoint) error {
	return u.repository.FindOrCreate(endpoint)
}

func (u endpointService) Update(endpoint *domain.Endpoint) error {
	return u.repository.Update(endpoint)
}

func (u endpointService) Delete(endpoint *domain.Endpoint) error {
	return u.repository.Delete(endpoint)
}
