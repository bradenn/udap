// Copyright (c) 2022 Braden Nicholson

package operators

import (
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
)

type triggerOperator struct {
	ctrl *controller.Controller
}

func NewTriggerOperator(ctrl *controller.Controller) ports.TriggerOperator {
	return &triggerOperator{
		ctrl: ctrl,
	}
}

func (m *triggerOperator) Run(trigger domain.Trigger) error {
	err := m.ctrl.SubRoutines.TriggerById(trigger.Id)
	if err != nil {
		return err
	}
	return nil
}
