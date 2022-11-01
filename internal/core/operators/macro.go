// Copyright (c) 2022 Braden Nicholson

package operators

import (
	"fmt"
	"time"
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
	"udap/internal/log"
	"udap/internal/pulse"
)

type macroOperator struct {
	ctrl  *controller.Controller
	holds map[string]chan bool
}

func NewMacroOperator(ctrl *controller.Controller) ports.MacroOperator {
	return &macroOperator{
		ctrl:  ctrl,
		holds: map[string]chan bool{},
	}
}

func (m *macroOperator) Run(macro domain.Macro) error {
	zone, err := m.ctrl.Zones.FindById(macro.ZoneId)
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

func (m *macroOperator) RunAndRevert(macro domain.Macro, baseline domain.Macro, revert time.Duration) error {
	zone, err := m.ctrl.Zones.FindById(macro.ZoneId)
	if err != nil {
		return err
	}

	values := map[string]string{}

	for _, entity := range zone.Entities {
		var attr *domain.Attribute
		attr, err = m.ctrl.Attributes.FindByComposite(entity.Id, macro.Type)
		if err != nil {
			return err
		}
		values[attr.Id] = attr.Value
		err = m.ctrl.Attributes.Request(entity.Id, macro.Type, macro.Value)
		if err != nil {
			continue
		}
	}

	if revert == 0 {
		return nil
	}

	go func() {
		begin := time.Now()
		tag := fmt.Sprintf("macro.%s.%s", macro.Id, begin)

		pulse.Begin(tag)
		defer pulse.End(tag)
		<-time.NewTimer(revert * time.Minute).C
		log.Event("REVERT TIME!!! %s (%s) (%dm)", tag, time.Since(begin), revert)
		zone, err = m.ctrl.Zones.FindById(baseline.ZoneId)
		if err != nil {
			return
		}
		for _, entity := range zone.Entities {
			var attr *domain.Attribute
			attr, err = m.ctrl.Attributes.FindByComposite(entity.Id, macro.Type)
			if err != nil {
				continue
			}
			err = m.ctrl.Attributes.Request(entity.Id, baseline.Type, values[attr.Id])
			if err != nil {
				continue
			}
		}

	}()

	return nil
}
