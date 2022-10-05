// Copyright (c) 2022 Braden Nicholson

package device

import (
	"time"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

type deviceService struct {
	repository  ports.DeviceRepository
	utilization map[string]domain.Utilization
	generic.Watchable[domain.Device]
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
	device.Utilization = u.utilization[device.Id]
	err := u.Emit(*device)
	if err != nil {
		return err
	}
	return nil
}

func (u *deviceService) Ping(id string, latency time.Duration) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}

	byId.LastSeen = time.Now()
	byId.Latency = latency
	byId.State = "ONLINE"

	err = u.repository.Update(byId)
	if err != nil {
		return err
	}
	err = u.emit(byId)
	if err != nil {
		return err
	}
	return nil
}

func (u *deviceService) Utilization(id string, utilization domain.Utilization) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	u.utilization[id] = utilization
	err = u.Update(byId)
	if err != nil {
		return err
	}
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

func NewService(repository ports.DeviceRepository) ports.DeviceService {
	return &deviceService{repository: repository, utilization: map[string]domain.Utilization{}}
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
