// Copyright (c) 2022 Braden Nicholson

package zone

import (
	"fmt"
	"udap/internal/core/domain"
)

type zoneService struct {
	repository domain.ZoneRepository
	channel    chan<- domain.Mutation
}

func (u *zoneService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, zone := range *all {
		err = u.emit(&zone)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *zoneService) emit(zone *domain.Zone) error {
	if u.channel == nil {
		return nil
	}
	u.channel <- domain.Mutation{
		Status:    "update",
		Operation: "zone",
		Body:      *zone,
		Id:        zone.Id,
	}
	return nil
}

func (u zoneService) FindByName(name string) (*domain.Zone, error) {
	return u.repository.FindByName(name)
}

func (u *zoneService) mutate(zone *domain.Zone) error {
	err := u.repository.Update(zone)
	if err != nil {
		return err
	}
	err = u.emit(zone)
	if err != nil {
		return err
	}
	return nil
}

func (u *zoneService) Watch(mut chan<- domain.Mutation) error {
	if u.channel != nil {
		return fmt.Errorf("channel already set")
	}
	u.channel = mut

	return nil
}

func NewService(repository domain.ZoneRepository) domain.ZoneService {
	return &zoneService{repository: repository}
}

func (u zoneService) Restore(id string) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}

	if !byId.Deleted {
		return fmt.Errorf("zone is already restored")
	}

	byId.Deleted = false

	err = u.mutate(byId)
	if err != nil {
		return err
	}

	return nil
}

func (u zoneService) Delete(id string) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}

	if byId.Deleted {
		return fmt.Errorf("zone is already deleted")
	}

	byId.Deleted = true

	err = u.mutate(byId)
	if err != nil {
		return err
	}

	return nil
}

// Repository Mapping

func (u zoneService) FindAll() (*[]domain.Zone, error) {
	return u.repository.FindAll()
}

func (u zoneService) FindById(id string) (*domain.Zone, error) {
	return u.repository.FindById(id)
}

func (u zoneService) Create(zone *domain.Zone) error {
	err := u.repository.Create(zone)
	if err != nil {
		return err
	}
	err = u.emit(zone)
	if err != nil {
		return err
	}
	return nil
}

func (u zoneService) FindOrCreate(zone *domain.Zone) error {
	return u.repository.FindOrCreate(zone)
}

func (u zoneService) Update(zone *domain.Zone) error {
	return u.repository.Update(zone)
}
