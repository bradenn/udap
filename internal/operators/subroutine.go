// Copyright (c) 2022 Braden Nicholson

package operators

import (
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
)

type subRoutineOperator struct {
	ctrl *controller.Controller
}

func NewSubRoutineOperator(ctrl *controller.Controller) ports.SubRoutineOperator {
	return &subRoutineOperator{
		ctrl: ctrl,
	}
}

func (m *subRoutineOperator) Run(subRoutine domain.SubRoutine) error {
	for _, macro := range subRoutine.Macros {
		// Good
		err := m.ctrl.Macros.Run(macro.Id)
		if err != nil {
			return err
		}
	}
	return nil
}
