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
	db     *gorm.DB
	router chi.Router
	server *http.Server

	controller *controller.Controller

	modules   domain.ModuleService
	endpoints domain.EndpointService
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

func (o *orchestrator) Run() error {

	resp := make(chan domain.Mutation, 20)
	o.controller.Listen(resp)
	wg := sync.WaitGroup{}

	wg.Add(3)
	go func() {
		defer wg.Done()
		for response := range resp {
			err := o.endpoints.SendAll(response.Id, response.Operation, response.Body)
			if err != nil {
				log.Err(err)
				continue
			}
		}
	}()

	// Initialize and route applicable domains
	routes.NewUserRouter(o.controller.Users).RouteUsers(o.router)
	routes.NewAttributeRouter(o.controller.Attributes).RouteAttributes(o.router)
	routes.NewEndpointRouter(o.endpoints).RouteEndpoints(o.router)
	routes.NewModuleRouter(o.modules).RouteModules(o.router)

	runtimes.NewModuleRuntime(o.modules)

	o.server = &http.Server{Addr: ":3020", Handler: o.router}

	go func() {
		defer wg.Done()
		err := o.server.ListenAndServe()
		if err != nil {
			log.ErrF(err, "http server exited with error:\n")
		}
		log.Event("Server exited")
	}()

	go func() {
		defer wg.Done()
		delay := time.Millisecond * 1000
		for {
			pulse.Begin("update")
			select {
			case <-time.After(delay):
				log.Event("Orchestrator event loop timed out")
				continue
			default:
				start := time.Now()

				err := o.Update()
				if err != nil {
					log.ErrF(err, "runtime update error: %s")
				}

				delta := time.Since(start)
				duration := delay - delta
				log.Event("Tick Complete (%s)", delta)
				if duration > 0 && duration < delay {
					time.Sleep(duration)
				}
				break
			}
			pulse.End("update")
			timings := pulse.Timings.Timings()
			for s, proc := range timings {
				resp <- domain.Mutation{
					Status:    "update",
					Operation: "timing",
					Body:      proc,
					Id:        s,
				}
			}
		}
		log.Event("Event loops exited")
	}()

	wg.Wait()
	return nil
}
