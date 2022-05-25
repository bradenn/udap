// Copyright (c) 2022 Braden Nicholson

package module

import (
	"udap/internal/core/domain"
)

type moduleService struct {
	repository domain.ModuleRepository
}

func NewModuleService(repository domain.ModuleRepository) domain.ModuleService {
	return moduleService{repository: repository}
}

// Repository Mapping

func (u moduleService) FindAll() ([]*domain.Module, error) {
	return u.repository.FindAll()
}

func (u moduleService) FindByName(name string) (*domain.Module, error) {
	return u.repository.FindByName(name)
}

func (u moduleService) FindById(id string) (*domain.Module, error) {
	// TODO implement me
	panic("implement me")
}

func (u moduleService) Disable(name string) error {
	// TODO implement me
	panic("implement me")
}

func (u moduleService) Enable(name string) error {
	// TODO implement me
	panic("implement me")
}

func (u moduleService) Reload(name string) error {
	// TODO implement me
	panic("implement me")
}

func (u moduleService) Halt(name string) error {
	// TODO implement me
	panic("implement me")
}
