// Copyright (c) 2022 Braden Nicholson

package orchestrator

import (
	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"net/http"
	"time"
	"udap/internal/controller"
	"udap/internal/core"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/modules/module"
	"udap/internal/port/routes"
	"udap/internal/port/runtimes"
	"udap/internal/pulse"
	"udap/platform/database"
	"udap/platform/router"
)

type orchestrator struct {
	db     *gorm.DB
	router chi.Router
	server *http.Server

	controller *controller.Controller

	modules domain.ModuleService
}

func (o *orchestrator) Update() error {
	endpoints, err := o.controller.Endpoints.FindAll()
	if err != nil {
		return err
	}
	eps := *endpoints
	for i := range eps {
		ep := eps[i]
		err = o.controller.Endpoints.Send(ep.Id, "endpoint", ep)
		if err != nil {
			return nil
		}
	}
	err = o.controller.Attributes.EmitAll()
	if err != nil {
		return err
	}
	return nil
}

func (o *orchestrator) Timings() error {
	timings := pulse.Timings.Timings()
	for _, timing := range timings {
		err := o.controller.Endpoints.Send("", "timing", timing)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *orchestrator) Run() error {

	o.server = &http.Server{Addr: ":3020", Handler: o.router}

	go func() {
		err := o.server.ListenAndServe()
		if err != nil {
			log.ErrF(err, "http server exited with error:\n")
		}
	}()

	attrs := make(chan domain.Attribute)
	err := o.controller.Attributes.Watch(attrs)
	if err != nil {
		return err
	}
	go func() {
		for attr := range attrs {
			err = o.controller.Endpoints.Send("", "attribute", attr)
			if err != nil {
				return
			}
		}
	}()

	delay := 1000.0
	for {
		select {
		case <-time.After(time.Millisecond * time.Duration(delay)):
			log.Event("Update timed out")
		default:
			start := time.Now()
			err := o.Update()
			if err != nil {
				log.ErrF(err, "runtime update error: %s")
			}
			d := time.Since(start)
			dur := (time.Millisecond * time.Duration(delay)) - d
			if dur > 0 {
				time.Sleep(dur)
			}
		}

		err := o.Timings()
		if err != nil {
			return err
		}

	}

	return nil
}

type Orchestrator interface {
	Start() error
	Run() error
}

func NewOrchestrator() Orchestrator {
	// Initialize Database
	db, err := database.New()
	if err != nil {
		return nil
	}
	// Initialize Router
	r := router.New()

	return &orchestrator{
		db:         db,
		router:     r,
		controller: nil,
	}
}

func (o *orchestrator) Start() error {
	var err error

	err = core.MigrateModels(o.db)
	if err != nil {
		return err
	}

	o.controller, err = controller.NewController(o.db)
	if err != nil {
		return err
	}

	o.modules = module.New(o.db, o.controller)

	// Initialize and route applicable domains
	routes.NewUserRouter(o.controller.Users).RouteUsers(o.router)
	routes.NewEndpointRouter(o.controller.Endpoints).RouteEndpoints(o.router)
	routes.NewModuleRouter(o.modules).RouteModules(o.router)

	runtimes.NewModuleRuntime(o.modules)

	return nil
}
