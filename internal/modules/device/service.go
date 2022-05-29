// Copyright (c) 2022 Braden Nicholson

package device

import (
	"udap/internal/core/domain"
)

type deviceService struct {
	repository domain.DeviceRepository
}

func (u deviceService) Register(device *domain.Device) error {
	return nil
}

func NewService(repository domain.DeviceRepository) domain.DeviceService {
	return deviceService{repository: repository}
}

// Repository Mapping

func (u deviceService) FindAll() ([]*domain.Device, error) {
	return u.repository.FindAll()
}

func (u deviceService) FindById(id string) (*domain.Device, error) {
	return u.repository.FindById(id)
}

func (u deviceService) Create(device *domain.Device) error {
	return u.repository.Create(device)
}

func (u deviceService) FindOrCreate(device *domain.Device) error {
	return u.repository.FindOrCreate(device)
}

func (u deviceService) Update(device *domain.Device) error {
	return u.repository.Update(device)
}

func (u deviceService) Delete(device *domain.Device) error {
	return u.repository.Delete(device)
}
