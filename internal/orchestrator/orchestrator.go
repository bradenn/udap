// Copyright (c) 2022 Braden Nicholson

package orchestrator

import (
	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"net/http"
	"sync"
	"time"
	"udap/internal/controller"
	"udap/internal/core"
	"udap/internal/core/domain"
	"udap/internal/core/modules/endpoint"
	"udap/internal/core/modules/module"
	"udap/internal/log"
	"udap/internal/port/routes"
	"udap/internal/port/runtimes"
	"udap/internal/pulse"
	"udap/platform/database"
	"udap/platform/router"
)

type orchestrator struct {
	db         *gorm.DB
	router     chi.Router
	server     *http.Server
	maxTick    time.Duration
	controller *controller.Controller

	modules   domain.ModuleService
	endpoints domain.EndpointService

	mutations chan domain.Mutation
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
		maxTick:    time.Second,
		mutations:  make(chan domain.Mutation, 16),
	}
}

func (o *orchestrator) Start() error {

	err := core.MigrateModels(o.db)
	if err != nil {
		return err
	}

	o.controller, err = controller.NewController(o.db)
	if err != nil {
		return err
	}

	o.modules = module.New(o.db, o.controller)
	o.endpoints = endpoint.New(o.db, o.controller)

	o.controller.Endpoints = o.endpoints
	o.controller.Modules = o.modules

	return nil
}

func (o *orchestrator) Update() error {
	err := o.modules.UpdateAll()
	if err != nil {
		return err
	}
	return nil
}

func (o *orchestrator) broadcastTimings() error {
	timings := pulse.Timings.Timings()
	for s, proc := range timings {
		o.mutations <- domain.Mutation{
			Status:    "update",
			Operation: "timing",
			Body:      proc,
			Id:        s,
		}
	}
	return nil
}

func (o *orchestrator) runServer() error {
	o.server = &http.Server{Addr: ":3020", Handler: o.router}
	err := o.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (o *orchestrator) handleMutations() error {
	for response := range o.mutations {
		err := o.endpoints.SendAll(response.Id, response.Operation, response.Body)
		if err != nil {
			log.Err(err)
			continue
		}
	}
	return nil
}

func (o *orchestrator) tick() <-chan error {
	out := make(chan error)
	go func() {
		start := time.Now()
		err := o.Update()
		if err != nil {
			out <- err
		}
		delta := time.Since(start)
		log.Event("Elapsed: %s", delta.String())
		if delta < o.maxTick && o.maxTick-delta > 250*time.Millisecond {
			time.Sleep(o.maxTick - delta - time.Millisecond*250)
		}
		out <- nil
	}()
	return out
}

func (o *orchestrator) Run() error {

	o.controller.WatchAll(o.mutations)

	wg := sync.WaitGroup{}

	wg.Add(3)
	go func() {
		defer wg.Done()
		err := o.handleMutations()
		if err != nil {
			log.Err(err)
			return
		}
	}()

	// Initialize and route applicable domains
	routes.NewUserRouter(o.controller.Users).RouteUsers(o.router)
	routes.NewAttributeRouter(o.controller.Attributes).RouteAttributes(o.router)
	routes.NewEndpointRouter(o.endpoints).RouteEndpoints(o.router)
	routes.NewModuleRouter(o.modules).RouteModules(o.router)

	runtimes.NewModuleRuntime(o.modules)

	go func() {
		defer wg.Done()
		err := o.runServer()
		if err != nil {
			log.Err(err)
			return
		}
	}()

	go func() {
		defer wg.Done()
		for {
			pulse.Begin("update")
			select {
			case <-time.After(o.maxTick + time.Millisecond*100):
				log.Event("Orchestrator event loop timed out")
				continue
			case err := <-o.tick():
				if err != nil {
					log.Err(err)
				}
			}
			pulse.End("update")
			err := o.broadcastTimings()
			if err != nil {
				log.Err(err)
			}
		}
	}()

	wg.Wait()
	return nil
}
