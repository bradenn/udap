// Copyright (c) 2022 Braden Nicholson

package zone

import (
	"fmt"
	"udap/internal/core/domain"
	"udap/internal/core/generic"
	"udap/internal/core/ports"
)

type zoneService struct {
	repository ports.ZoneRepository
	generic.Watchable[domain.Zone]
}

func (u *zoneService) EmitAll() error {
	all, err := u.FindAll()
	if err != nil {
		return err
	}
	for _, zone := range *all {
		err = u.Emit(zone)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *zoneService) FindByName(name string) (*domain.Zone, error) {
	return u.repository.FindByName(name)
}

func (u *zoneService) mutate(zone *domain.Zone) error {
	err := u.repository.Update(zone)
	if err != nil {
		return err
	}
	err = u.Emit(*zone)
	if err != nil {
		return err
	}
	return nil
}

func NewService(repository ports.ZoneRepository) ports.ZoneService {
	return &zoneService{repository: repository}
}

func (u *zoneService) Restore(id string) error {
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

func (u *zoneService) Delete(id string) error {
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

func (u *zoneService) AddEntity(id string, entity string) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	e := domain.Entity{}
	e.Id = entity
	byId.Entities = append(byId.Entities, e)
	err = u.mutate(byId)
	if err != nil {
		return err
	}
	return nil
}

func (u *zoneService) RemoveEntity(id string, entity string) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	var entities []domain.Entity

	for _, e2 := range byId.Entities {
		if e2.Id != entity {
			entities = append(entities, e2)
		}
	}
	byId.Entities = entities
	err = u.mutate(byId)
	if err != nil {
		return err
	}
	return nil
}

func (u *zoneService) Pin(id string) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	byId.Pinned = true
	err = u.mutate(byId)
	if err != nil {
		return err
	}
	return nil
}

func (u *zoneService) Unpin(id string) error {
	byId, err := u.repository.FindById(id)
	if err != nil {
		return err
	}
	byId.Pinned = false
	err = u.mutate(byId)
	if err != nil {
		return err
	}
	return nil
}

// Repository Mapping

func (u *zoneService) FindAll() (*[]domain.Zone, error) {
	return u.repository.FindAll()
}

func (u *zoneService) FindById(id string) (*domain.Zone, error) {
	return u.repository.FindById(id)
}

func (u *zoneService) Create(zone *domain.Zone) error {
	err := u.repository.Create(zone)
	if err != nil {
		return err
	}
	err = u.Emit(*zone)
	if err != nil {
		return err
	}
	return nil
}

func (u *zoneService) FindOrCreate(zone *domain.Zone) error {
	return u.repository.FindOrCreate(zone)
}

func (u *zoneService) Update(zone *domain.Zone) error {
	return u.repository.Update(zone)
}
