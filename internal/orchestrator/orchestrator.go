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
	"udap/internal/core/device"
	"udap/internal/core/domain"
	"udap/internal/log"
	"udap/internal/modules"
	"udap/internal/pulse"
	"udap/internal/srv"
	"udap/internal/srv/store"
	"udap/platform/database"
)

type orchestrator struct {
	db         *gorm.DB
	controller *controller.Controller
	store      *store.Store
	server     srv.Server
	maxTick    time.Duration
	done       chan bool
	ready      bool
	sys        srv.System
	mutations  chan domain.Mutation
}

type Orchestrator interface {
	Start() error
	Run() error
}

func (o *orchestrator) Terminate() {

	go func() {
		err := o.controller.Endpoints.CloseAll()
		if err != nil {
			log.Err(err)
		}
	}()

	go func() {
		err := o.controller.Modules.DisposeAll()
		if err != nil {
			log.Err(err)
		}
	}()

	close(o.mutations)

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
	server := srv.NewServer()

	str := store.NewStore()

	// Initialize Orchestrator
	return &orchestrator{
		db:         db,
		server:     server,
		store:      str,
		done:       make(chan bool),
		controller: nil,
		maxTick:    time.Second * 1,
		mutations:  make(chan domain.Mutation, 512),
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

	go func() {
		err := o.handleMutations()
		if err != nil {
			log.Err(err)
			return
		}
	}()

	err := core.MigrateModels(o.db)
	if err != nil {
		return err
	}

	o.controller, err = controller.NewController(o.mutations)
	if err != nil {
		return err
	}

	o.ready = false

	o.sys = srv.NewRtx(&o.server, o.controller, o.db, o.store)

	o.sys.UseModules(
		modules.NewModule, modules.NewTrace)

	o.sys.UseModules(
		modules.NewEntity,
		modules.NewAttribute,
		modules.NewZone)

	o.sys.UseModules(
		modules.NewEndpoint)

	o.sys.UseModules(
		modules.NewMacro,
		modules.NewSubroutine,
		modules.NewTrigger,
		modules.NewUser,
		modules.NewNetwork,
		device.NewModule,
		modules.NewNotifications,
		modules.NewLog,
	)

	o.sys.UseModules(modules.NewAction)

	o.sys.Loaded()
	o.ready = true

	return nil
}

func (o *orchestrator) Update() error {
	if !o.ready {
		return nil
	}
	//err := o.controller.Modules.UpdateAll()
	//if err != nil {
	//	return err
	//}
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
		for !o.ready {
			time.Sleep(time.Millisecond * 250)
		}

		err := o.controller.Endpoints.SendAll(response.Id, response.Operation, response.Body)
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
		pulse.End("update")
		_ = o.broadcastTimings()
		if err != nil {
			out <- err
			return
		}

		delta := time.Since(start)
		if delta < o.maxTick {
			// log.Tick("Elapsed: %s", delta.String())
			time.Sleep(o.maxTick - delta)
		}
		select {
		case out <- nil:
			break
		default:
			break
		}
	}()
	return out
}

func (o *orchestrator) Run() error {

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		err := o.server.Run()
		if err != nil {
			log.Err(err)
			return
		}
	}()

	t := time.NewTimer(o.maxTick + time.Millisecond*500)

	go func() {
		defer wg.Done()

		for {
			pulse.Begin("update")
			t.Reset(o.maxTick + time.Second*500)
			select {
			case <-o.done:
				log.Event("Event loop exiting...")
				o.Terminate()
				return
			case <-t.C:
				log.Event("Orchestrator event loop timed out (%s)", (o.maxTick + time.Millisecond*500).String())
				log.Event("Currently %d threads.", runtime.NumGoroutine())

			case err := <-o.tick():
				t.Stop()
				if err != nil {
					log.Err(err)
				}
			}

		}
	}()

	wg.Wait()
	return nil
}
