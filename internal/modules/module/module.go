// Copyright (c) 2022 Braden Nicholson

package module

import (
	"udap/internal/core/domain"
)

func New() domain.ModuleService {
	repo := NewRepository()
	service := NewService(repo)
	return service
}
