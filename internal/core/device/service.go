// Copyright (c) 2022 Braden Nicholson

package device

import (
	"time"
	"udap/internal/core/generic"
)

func newService(repository Repository) Service {
	return &deviceService{repository: repository, utilization: map[string]Utilization{}}
}

type deviceService struct {
	repository  Repository
	utilization map[string]Utilization
	generic.Watchable[Device]
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

func (u *deviceService) emit(device *Device) error {
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

func (u *deviceService) Utilization(id string, utilization Utilization) error {
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

func (u *deviceService) Register(device *Device) error {
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

// Repository Mapping

func (u *deviceService) FindAll() (*[]Device, error) {
	return u.repository.FindAll()
}

func (u *deviceService) FindById(id string) (*Device, error) {
	return u.repository.FindById(id)
}

func (u *deviceService) Create(device *Device) error {
	return u.repository.Create(device)
}

func (u *deviceService) FindOrCreate(device *Device) error {
	return u.repository.FindOrCreate(device)
}

func (u *deviceService) Update(device *Device) error {

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

func (u *deviceService) Delete(device *Device) error {
	return u.repository.Delete(device)
}
