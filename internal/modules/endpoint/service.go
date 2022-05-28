// Copyright (c) 2022 Braden Nicholson

package endpoint

import (
	"udap/internal/core/domain"
)

type endpointService struct {
	repository domain.EndpointRepository
}

func (u endpointService) FindByKey(key string) (*domain.Endpoint, error) {
	return u.repository.FindByKey(key)
}

func (u endpointService) Enroll(key string) (*domain.Endpoint, error) {
	// TODO implement me
	panic("implement me")
}

func (u endpointService) Disconnect(key string) error {
	// TODO implement me
	panic("implement me")
}

func NewService(repository domain.EndpointRepository) domain.EndpointService {
	return endpointService{repository: repository}
}

// Repository Mapping

func (u endpointService) FindAll() ([]*domain.Endpoint, error) {
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
