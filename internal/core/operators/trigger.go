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

func (m *triggerOperator) RunCustom(trigger domain.Trigger, key string, value string) error {

	actions, err := m.ctrl.Actions.FindByTriggerId(trigger.Id)
	if err != nil {
		return err
	}

	for _, action := range *actions {
		err = m.ctrl.Actions.ExecuteCustomById(action.Id, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *triggerOperator) Run(trigger domain.Trigger) error {

	err := m.ctrl.SubRoutines.TriggerById(trigger.Id)
	if err != nil {
		return err
	}

	actions, err := m.ctrl.Actions.FindByTriggerId(trigger.Id)
	if err != nil {
		return err
	}

	for _, action := range *actions {
		err = m.ctrl.Actions.ExecuteById(action.Id)
		if err != nil {
			return err
		}
	}

	return nil
}
