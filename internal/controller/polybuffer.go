// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"udap/internal/log"
)

type Observer struct {
	caller   string
	callback func(data interface{}) error
}

type Mutation struct {
	Key   string
	Value interface{}
}

type Observable struct {
	watchers       map[string]Observer
	singleWatchers map[string]map[string]Observer
	handler        chan Mutation
}

func (p *Observable) Run() {
	p.watchers = map[string]Observer{}
	p.singleWatchers = map[string]map[string]Observer{}
	p.handler = make(chan Mutation, 1)
	go func() {
		// Individual-key observers
		for a := range p.handler {
			if p.singleWatchers[a.Key] != nil {
				for _, watcher := range p.singleWatchers[a.Key] {
					err := watcher.callback(a.Value)
					if err != nil {
						log.Err(err)
					}

				}
			}
			// All-key observers
			for _, watcher := range p.watchers {
				err := watcher.callback(a.Value)
				if err != nil {
					log.Err(err)
				}
			}
		}
	}()
}

func (p *Observable) emit(id string, data interface{}) {
	p.handler <- Mutation{
		Key:   id,
		Value: data,
	}
}

func (p *Observable) WatchSingle(key string, fn func(data interface{}) error) {
	_, file, _, ok := runtime.Caller(1)
	if ok {
		path := filepath.Base(file)
		o := Observer{
			caller:   path,
			callback: fn,
		}
		if p.singleWatchers[key] == nil {
			p.singleWatchers[key] = map[string]Observer{}
		}

		if p.singleWatchers[key][path].caller != "" {
			log.Err(fmt.Errorf("already watching"))
		}
		p.singleWatchers[key][path] = o
	}
}

func (p *Observable) Watch(fn func(data interface{}) error) {
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
	raw  map[string]interface{}
	data sync.Map
}

func (p *PolyBuffer) set(id string, data interface{}) {
	p.data.Store(id, data)
}

func (p *PolyBuffer) get(name string) interface{} {
	res, ok := p.data.Load(name)
	if !ok {
		return nil
	}
	return res
}

func (p *PolyBuffer) Keys() []string {
	var s []string
	p.data.Range(func(key, value interface{}) bool {
		s = append(s, key.(string))
		return true
	})
	return s
}
