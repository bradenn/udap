// Copyright (c) 2022 Braden Nicholson

package macro

import (
	"udap/internal/controller"
	"udap/internal/core/domain"
)

type macroOperator struct {
	ctrl *controller.Controller
}

func NewOperator(ctrl *controller.Controller) domain.MacroOperator {
	return &macroOperator{
		ctrl: ctrl,
	}
}

func (m *macroOperator) Run(macro domain.Macro) error {
	zone, err := m.ctrl.Zones.FindById(macro.Id)
	if err != nil {
		return err
	}

	for _, entity := range zone.Entities {
		err = m.ctrl.Attributes.Request(entity.Id, macro.Type, macro.Value)
		if err != nil {
			continue
		}
	}
	return nil
}
