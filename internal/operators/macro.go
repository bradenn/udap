// Copyright (c) 2022 Braden Nicholson

package operators

import (
	"fmt"
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
)

type macroOperator struct {
	ctrl *controller.Controller
}

func NewMacroOperator(ctrl *controller.Controller) ports.MacroOperator {
	return &macroOperator{
		ctrl: ctrl,
	}
}

func (m *macroOperator) Run(macro domain.Macro) error {
	zone, err := m.ctrl.Zones.FindById(macro.ZoneId)
	if err != nil {
		return err
	}

	fmt.Println(zone)

	fmt.Printf("Running macro: '%s', Zone: '%s' @ %d\n", macro.Name, zone.Name, len(zone.Entities))

	for _, entity := range zone.Entities {
		fmt.Printf("MACRO RUN %s\n", entity.Name)
		err = m.ctrl.Attributes.Request(entity.Id, macro.Type, macro.Value)
		if err != nil {
			continue
		}
	}
	return nil
}
