// Copyright (c) 2021 Braden Nicholson

package controllers

import "udap/internal/models"

type EntityController struct {
	entity *models.Entity
}

func NewEntityController(entity *models.Entity) EntityController {
	return EntityController{
		entity: entity,
	}
}

func (e *EntityController) PushState(state models.State) error {
	err := e.entity.Push(state)
	if err != nil {
		return err
	}
	return nil
}

func (e *EntityController) GetState() (models.State, error) {
	err := e.entity.Poll()
	if err != nil {
		return nil, err
	}
	return e.entity.State, nil
}
