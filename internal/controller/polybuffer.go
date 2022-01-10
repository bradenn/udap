// Copyright (c) 2021 Braden Nicholson

package controller

import "sync"

type PolyBuffer struct {
	raw  map[string]any
	data sync.Map
}

func (p *PolyBuffer) load() {
	p.data = sync.Map{}
}

func (p *PolyBuffer) set(id string, data any) {
	if p.raw == nil {
		p.raw = map[string]any{}
	}
	p.data.Store(id, data)
}

func (p *PolyBuffer) peek(name string) any {
	res, ok := p.data.Load(name)
	if !ok {
		return nil
	}
	return res
}

func (p *PolyBuffer) get(name string) any {
	res, ok := p.data.Load(name)
	if !ok {
		return nil
	}
	return res
}

func (p *PolyBuffer) Keys() []string {
	var s []string
	p.data.Range(func(key, value any) bool {
		s = append(s, key.(string))
		return true
	})
	// for s2, _ := range p.raw {
	// 	s = append(s, s2)
	// }
	return s
}
