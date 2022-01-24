// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"sync"
)

type PolyBuffer struct {
	raw  map[string]any
	data sync.Map
}

func (p *PolyBuffer) set(id string, data any) {
	p.data.Store(id, data)
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
	return s
}
