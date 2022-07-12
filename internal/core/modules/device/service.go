// Copyright (c) 2022 Braden Nicholson

package device

import (
	"fmt"
	"udap/internal/core/domain"
)

type deviceService struct {
	repository domain.DeviceRepository
	channel    chan<- domain.Mutation
}

func (u *deviceService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, device := range *all {
		err = u.emit(&device)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *deviceService) emit(device *domain.Device) error {
	if u.channel == nil {
		return nil
	}
	u.channel <- domain.Mutation{
		Status:    "update",
		Operation: "device",
		Body:      *device,
		Id:        device.Id,
	}
	return nil
}

func (u *deviceService) Watch(mut chan<- domain.Mutation) error {
	if u.channel != nil {
		return fmt.Errorf("channel already set")
	}
	u.channel = mut

	return nil
}

func (u deviceService) Register(device *domain.Device) error {
	err := u.repository.FindOrCreate(device)
	if err != nil {
		return err
	}
	err = u.emit(device)
	if err != nil {
		return err
	}
	return nil
}

func NewService(repository domain.DeviceRepository) domain.DeviceService {
	return &deviceService{repository: repository}
}

// Repository Mapping

func (u deviceService) FindAll() (*[]domain.Device, error) {
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

	err := u.repository.Update(device)
	if err != nil {
		return err
	}

	err = u.emit(device)
	if err != nil {
		return err
	}

	return nil
}

func (u deviceService) Delete(device *domain.Device) error {
	return u.repository.Delete(device)
}
