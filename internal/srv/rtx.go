// Copyright (c) 2022 Braden Nicholson

package srv

import (
	"gorm.io/gorm"
	"udap/internal/controller"
	"udap/internal/core/domain"
	"udap/internal/port/routes"
	"udap/internal/srv/store"
)

type Watch interface {
	Watch(mutation chan<- domain.Mutation)
	EmitAll() error
}

type Context interface {
}

type CoreModule interface {
	Mount(rtx Context) error
	Unmount() error
	Mounted() bool
}

type System interface {
	WithWatch(watch Watch)
	WhenLoaded(watch func())
	WithRoute(route routes.Routable)
	UseModules(modules ...func(sys System))
	DB() *gorm.DB
	Store() *store.Store
	Loaded()
	Ctrl() *controller.Controller
}

type sys struct {
	db    *gorm.DB
	store *store.Store
	*Server
	ctrl   *controller.Controller
	onLoad chan bool
	loaded bool
}

func (r *sys) WithWatch(mutation Watch) {

	mutation.Watch(r.ctrl.RX)

	err := mutation.EmitAll()
	if err != nil {
		return
	}
}

func (r *sys) WhenLoaded(watch func()) {
	go func() {
		select {
		case load := <-r.onLoad:
			if load {
				watch()
			}
		}
	}()
}

func (r *sys) Loaded() {
	r.onLoad <- true
}

func (r *sys) DB() *gorm.DB {
	return r.db
}

func (r *sys) Store() *store.Store {
	return r.store
}

func (r *sys) Ctrl() *controller.Controller {
	return r.ctrl
}

func (r *sys) WithRoute(route routes.Routable) {
	r.AddRoute(route)
}

func (r *sys) UseModules(modules ...func(sys System)) {
	for _, module := range modules {
		module(r)
	}
}

func NewRtx(server *Server, ctrl *controller.Controller, db *gorm.DB, str *store.Store) System {
	return &sys{
		db:     db,
		Server: server,
		ctrl:   ctrl,
		store:  str,
		onLoad: make(chan bool),
	}
}
