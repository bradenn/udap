// Copyright (c) 2022 Braden Nicholson

package zone

import (
	"udap/internal/core/domain"
)

type zoneService struct {
	repository domain.ZoneRepository
}

func NewService(repository domain.ZoneRepository) domain.ZoneService {
	return zoneService{repository: repository}
}

// Repository Mapping

func (u zoneService) FindAll() ([]*domain.Zone, error) {
	return u.repository.FindAll()
}

func (u zoneService) FindById(id string) (*domain.Zone, error) {
	return u.repository.FindById(id)
}

func (u zoneService) Create(zone *domain.Zone) error {
	return u.repository.Create(zone)
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
