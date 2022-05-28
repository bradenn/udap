// Copyright (c) 2022 Braden Nicholson

package module

import (
	"udap/internal/core/domain"
)

type moduleRepo struct {
}

func NewRepository() domain.ModuleRepository {
	return &moduleRepo{}
}

func (u moduleRepo) Candidates() ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (u moduleRepo) FindAll() ([]*domain.Module, error) {
	// TODO implement me
	panic("implement me")
}

func (u moduleRepo) FindByName(name string) (*domain.Module, error) {
	// TODO implement me
	panic("implement me")
}
