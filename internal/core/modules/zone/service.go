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

func (u zoneService) Delete(zone *domain.Zone) error {
	return u.repository.Delete(zone)
}
