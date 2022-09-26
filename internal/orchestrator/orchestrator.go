// Copyright (c) 2022 Braden Nicholson

package orchestrator

import (
	"fmt"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
	"udap/internal/controller"
	"udap/internal/core"
	"udap/internal/core/domain"
	"udap/internal/core/services/endpoint"
	"udap/internal/core/services/module"
	"udap/internal/log"
	"udap/internal/port/routes"
	"udap/internal/port/runtimes"
	"udap/internal/pulse"
	"udap/platform/database"
)

type orchestrator struct {
	db         *gorm.DB
	controller *controller.Controller

	server  Server
	maxTick time.Duration
	done    chan bool

	modules   domain.ModuleService
	endpoints domain.EndpointService

	mutations chan domain.Mutation
}

type Orchestrator interface {
	Start() error
	Run() error
}

func (o *orchestrator) Terminate(reason string) {
	_ = o.modules.DisposeAll()
	_ = o.endpoints.CloseAll()
	fmt.Printf("\nThreads at exit: %d\n", runtime.NumGoroutine())
	os.Exit(0)
}

func NewOrchestrator() (Orchestrator, error) {
	// Initialize Database
	db, err := database.New()
	if err != nil {
		return nil, err
	}
	// Initialize Server
	server := NewServer()
	// Initialize Orchestrator
	return &orchestrator{
		db:         db,
		server:     server,
		done:       make(chan bool),
		controller: nil,
		maxTick:    time.Second,
		mutations:  make(chan domain.Mutation, 32),
	}, nil
}

func (o *orchestrator) Start() error {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			_ = o.server.Close()
			o.done <- true
			return
		}
	}()

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

func (o *orchestrator) handleMutations() error {
	for response := range o.mutations {
		// o.modules.HandleEmits(response)
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
			return
		}
		delta := time.Since(start)
		if delta < o.maxTick {
			// log.Tick("Elapsed: %s", delta.String())
			time.Sleep(o.maxTick - delta)
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
	o.server.AddRoutes(
		routes.NewEndpointRouter(o.controller.Endpoints),
		routes.NewAttributeRouter(o.controller.Attributes),
		routes.NewUserRouter(o.controller.Users),
		routes.NewZoneRouter(o.controller.Zones),
		routes.NewDeviceRouter(o.controller.Devices),
		routes.NewEntityRouter(o.controller.Entities),
		routes.NewModuleRouter(o.modules),
		routes.NewEndpointRouter(o.endpoints),
	)

	// o.router.Group(func(r chi.Router) {
	//
	// 	routes.NewUserRouter(o.controller.Users).RouteInternal(r)
	// 	routes.NewAttributeRouter(o.controller.Attributes).RouteInternal(r)
	// 	routes.NewZoneRouter(o.controller.Zones).RouteInternal(r)
	// 	routes.NewDeviceRouter(o.controller.Devices).RouteInternal(r)
	// 	routes.NewEntityRouter(o.controller.Entities).RouteInternal(r)
	// 	routes.NewModuleRouter(o.modules).RouteInternal(r)
	// 	routes.NewEndpointRouter(o.endpoints).RouteInternal(r)
	// })
	// routes.NewEndpointRouter(o.endpoints).RouteExternal(o.router)
	runtimes.NewModuleRuntime(o.modules)

	go func() {
		defer wg.Done()
		err := o.server.Run()
		if err != nil {
			// log.Err(err)
			return
		}
	}()

	t := time.NewTimer(o.maxTick + time.Millisecond*100)

	go func() {
		defer wg.Done()
		for {
			pulse.Begin("update")
			t.Reset(o.maxTick + time.Millisecond*100)
			select {
			case <-o.done:
				log.Event("Event loop exiting...")
				o.Terminate("Terminated")
				close(o.mutations)
				return
			case <-t.C:
				log.Event("Orchestrator event loop timed out")
				pulse.End("update")
				continue
			case err := <-o.tick():
				t.Stop()
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
