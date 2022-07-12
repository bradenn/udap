// Copyright (c) 2022 Braden Nicholson

package runtimes

import (
	"udap/internal/core/domain"
	"udap/internal/log"
)

func NewModuleRuntime(service domain.ModuleService) {
	err := service.Discover()
	if err != nil {
		return
	}

	err = service.BuildAll()
	if err != nil {
		log.Err(err)
		return
	}
	err = service.LoadAll()
	if err != nil {
		return
	}

	err = service.RunAll()
	if err != nil {
		return
	}

}
