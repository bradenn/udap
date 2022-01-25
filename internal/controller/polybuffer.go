// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"path/filepath"
	"runtime"
	"sync"
)

type Observer struct {
	caller   string
	callback func(data any) error
}

type Mutation struct {
	Key   string
	Value any
}

type Observable struct {
	watchers map[string]Observer
	handler  chan Mutation
}

func (p *Observable) Run() {
	p.watchers = map[string]Observer{}
	p.handler = make(chan Mutation)
	go func() {
		for a := range p.handler {
			for _, watcher := range p.watchers {
				err := watcher.callback(a.Value)
				if err != nil {
					return
				}
			}
		}
	}()
}

func (p *Observable) emit(id string, data any) {
	p.handler <- Mutation{
		Key:   id,
		Value: data,
	}
}

func (p *Observable) Watch(fn func(data any) error) {
	_, file, _, ok := runtime.Caller(1)
	if ok {
		path := filepath.Base(file)
		o := Observer{
			caller:   path,
			callback: fn,
		}
		p.watchers[path] = o
	}

}

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
