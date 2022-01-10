// Copyright (c) 2021 Braden Nicholson

package server

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"udap/internal/bond"
	"udap/internal/controller"
	"udap/internal/log"
)

type Daemon interface {
	Setup(bond *bond.Bond) error
	Run() error
	Update() error
}

// Runtime manages the event-loop for all instances, as well as the websocket connections between UDAP and endpoints.
type Runtime struct {
	ctrl         *controller.Controller
	daemons      []Daemon
	eventHandler chan bond.Msg

	Endpoints *Endpoints
	Modules   *Modules
}

// Dependency is the level at which this service needs to run
func (r *Runtime) Dependency() (level int) {
	return 1
}

func (r *Runtime) Channel() chan bond.Msg {
	return r.eventHandler
}

func (r *Runtime) Name() (name string) {
	return "runtime"
}

type RuntimeStats struct {
	Threads int `json:"threads"`
}

func (r *Runtime) logRuntimeData() {
	_ = RuntimeStats{
		Threads: runtime.NumGoroutine(),
	}
}

func (r *Runtime) handleRequest() {
	for msg := range r.eventHandler {
		// start := time.Now()
		msg.Respond(r.ctrl.Handle(msg))
		// log.Event("EVENT: %s.%s (%s)", msg.Target, msg.Operation, time.Since(start))
	}
	close(r.eventHandler)
}

func (r *Runtime) Update() error {
	err := r.UpdateDaemons()
	if err != nil {
		return err
	}
	return nil
}

func (r *Runtime) AddDaemons(daemon ...Daemon) {
	for _, d := range daemon {
		r.daemons = append(r.daemons, d)
	}
}

func (r *Runtime) SetupDaemons() (err error) {
	b := bond.NewBond(r.eventHandler)
	wg := sync.WaitGroup{}
	wg.Add(len(r.daemons))
	for i, d := range r.daemons {
		go func(daemon Daemon, id int) {
			defer wg.Done()
			log.Log("Daemon '%d' loaded.", id)
			err = daemon.Setup(b)
			if err != nil {
				return
			}
		}(d, i)
	}
	wg.Wait()
	return nil
}

func (r *Runtime) RunDaemons() (err error) {

	return nil
}

func (r *Runtime) UpdateDaemons() (err error) {
	wg := sync.WaitGroup{}
	wg.Add(len(r.daemons))
	for _, d := range r.daemons {
		go func(daemon Daemon) {
			defer wg.Done()
			err = daemon.Update()
			if err != nil {
				log.Err(err)
				return
			}
		}(d)
	}
	wg.Wait()
	return nil
}

func (r *Runtime) Load() (err error) {
	r.daemons = []Daemon{}

	r.Modules = &Modules{}
	r.Endpoints = &Endpoints{}

	r.AddDaemons(r.Modules, r.Endpoints)

	r.ctrl, err = controller.NewController()
	if err != nil {
		return err
	}

	r.eventHandler = make(chan bond.Msg, 16)

	err = r.SetupDaemons()
	if err != nil {
		return err
	}

	return nil
}

// Run is called when the runtime is to begin accepting traffic
func (r *Runtime) Run() (err error) {
	go r.handleRequest()
	wg := sync.WaitGroup{}
	wg.Add(len(r.daemons))
	for i, d := range r.daemons {
		log.Log("Daemon '%d' running.", i)
		go func(daemon Daemon) {
			defer wg.Done()
			err = daemon.Run()
			if err != nil {
				log.Err(err)
				return
			}
		}(d)
	}

	go func() {
		for {
			start := time.Now()
			err = r.Update()
			if err != nil {
				log.ErrF(err, "runtime update error: %s")
			}
			d := time.Since(start)
			log.Event("Update: %d threads, %s", runtime.NumGoroutine(), d.String())
			select {
			case <-time.After(time.Millisecond * 3000):
				log.ErrF(fmt.Errorf("timed out main update loop"), "%s")
				continue
			default:
				time.Sleep(time.Millisecond*3000 - d)
			}
		}
	}()

	wg.Wait()

	return nil
}
