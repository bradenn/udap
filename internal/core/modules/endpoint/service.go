// Copyright (c) 2022 Braden Nicholson

package endpoint

import (
	"fmt"
	"github.com/gorilla/websocket"
	"udap/internal/core/domain"
)

type endpointService struct {
	repository domain.EndpointRepository
	operator   domain.EndpointOperator
	channel    chan<- domain.Mutation
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
		err = u.emit(&endpoint)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *endpointService) emit(endpoint *domain.Endpoint) error {
	if u.channel == nil {
		return fmt.Errorf("channel is null")
	}
	u.channel <- domain.Mutation{
		Status:    "update",
		Operation: "endpoint",
		Body:      *endpoint,
		Id:        endpoint.Id,
	}
	return nil
}

func (u *endpointService) Watch(ref chan<- domain.Mutation) error {
	if u.channel != nil {
		return fmt.Errorf("channel in use")
	}
	u.channel = ref
	return nil
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

func (u *endpointService) Enroll(id string, conn *websocket.Conn) error {
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

func (u *endpointService) Disconnect(key string) error {
	err := u.operator.Unenroll(key)
	if err != nil {
		return err
	}
	return nil
}

func NewService(repository domain.EndpointRepository, operator domain.EndpointOperator) domain.EndpointService {
	return &endpointService{repository: repository, operator: operator}
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
