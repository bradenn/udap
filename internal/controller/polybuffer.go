// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"fmt"
	"sync"
)

type Gen[T any] struct {
	data sync.Map
}

func (p *Gen[T]) Keys() []string {
	var s []string
	p.data.Range(func(key, value any) bool {
		s = append(s, key.(string))
		return true
	})
	return s
}

func (p *Gen[T]) set(id string, data T) error {
	p.data.Store(id, data)
	return nil
}

func (p *Gen[T]) get(id string) (t T, err error) {
	val, ok := p.data.Load(id)
	if !ok {
		return t, fmt.Errorf("mutable '%s' not found", id)
	}
	t = val.(T)
	return t, nil
}

type PolyBuffer struct {
	raw  map[string]any
	data sync.Map
}

func (p *PolyBuffer) Put(id string, data any) error {
	p.data.Store(id, data)
	return nil
}

func (p *PolyBuffer) Get(id string) (any, error) {
	res, ok := p.data.Load(id)
	if !ok {
		return res, fmt.Errorf("item does not exist at this address")
	}
	return res, nil
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
