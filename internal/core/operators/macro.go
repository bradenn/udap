// Copyright (c) 2022 Braden Nicholson

package operators

import (
	"time"
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
)

type macroOperator struct {
	ctrl     *controller.Controller
	launcher chan<- domain.Macro
	holds    map[string]*time.Timer
	revert   map[string]domain.Macro
}

func NewMacroOperator(ctrl *controller.Controller) ports.MacroOperator {
	mo := &macroOperator{
		ctrl:   ctrl,
		holds:  map[string]*time.Timer{},
		revert: map[string]domain.Macro{},
	}

	//go func() {
	//	select {
	//	case <-launcher:
	//
	//	}
	//
	//	for macro, timer := range mo.holds {
	//		mo.runRevert(macro)
	//	}
	//}()

	return mo
}

func (m *macroOperator) runRevert(macro domain.Macro) error {
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
	//if m.holds[macro.Id] != nil {
	//	m.holds[macro.Id].Stop()
	//	delete(m.holds, macro.Id)
	//	delete(m.revert, macro.Id)
	//}

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

	////begin := time.Now()
	////tag := fmt.Sprintf("macro.%s.%s", macro.Id, begin)
	//if m.holds[macro.Id] != nil {
	//	m.holds[macro.Id].Reset(revert * time.Minute)
	//} else {
	//	m.holds[macro.Id] = time.NewTimer(revert * time.Minute)
	//	go func() {
	//		_ = <-m.holds[macro.Id].C
	//		if m.holds[macro.Id] != nil {
	//			log.Event("REVERTING '%s'.", macro.Name)
	//			err = m.runRevert(m.revert[macro.Id])
	//			if err != nil {
	//				return
	//			}
	//			delete(m.holds, macro.Id)
	//			delete(m.revert, macro.Id)
	//		}
	//
	//	}()
	//}

	return nil
}
