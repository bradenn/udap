// Copyright (c) 2021 Braden Nicholson

package server

import (
	"sync"
	"time"
	"udap/internal/bond"
	"udap/internal/controller"
	"udap/internal/log"
	"udap/internal/pulse"
)

type Daemon interface {
	Setup(ctrl *controller.Controller, bond *bond.Bond) error
	Name() string
	Run() error
	Update() error
}

// Runtime manages the event-loop for all instances, as well as the websocket connections between UDAP and endpoints.
type Runtime struct {
	ctrl         *controller.Controller
	daemons      []Daemon
	eventHandler chan bond.Msg

	System System

	Endpoints *Endpoints
	Modules   *Modules
}

func (r *Runtime) Name() (name string) {
	return "runtime"
}

func (r *Runtime) handleRequest() {
	for msg := range r.eventHandler {
		start := time.Now()
		msg.Respond(r.ctrl.Handle(msg))
		log.Event("EVENT: %s.%s (%s)", msg.Target, msg.Operation, time.Since(start))
	}
}

func (r *Runtime) Update() error {
	wg := sync.WaitGroup{}
	wg.Add(len(r.daemons))
	for _, d := range r.daemons {
		go func(daemon Daemon) {
			defer wg.Done()
			err := daemon.Update()
			if err != nil {
				log.Err(err)
			}
		}(d)
	}
	wg.Wait()
	return nil
}

func (r *Runtime) addDaemons(daemon ...Daemon) {
	for _, d := range daemon {
		r.daemons = append(r.daemons, d)
	}
}

func (r *Runtime) SetupDaemons() (err error) {
	b := bond.NewBond(r.eventHandler)
	for _, d := range r.daemons {
		log.Log("Daemon '%s' loaded.", d.Name())
		err = d.Setup(r.ctrl, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Runtime) Load() (err error) {

	r.daemons = []Daemon{}

	r.Modules = &Modules{}
	r.Endpoints = &Endpoints{}

	r.addDaemons(r.Modules, r.Endpoints)

	r.eventHandler = make(chan bond.Msg, 8)

	r.ctrl, err = controller.NewController()
	if err != nil {
		return err
	}

	info, err := systemInfo()
	if err != nil {
		return err
	}
	SystemInfo = info

	err = r.SetupDaemons()
	if err != nil {
		return err
	}

	return nil
}

// Run is called when the runtime is to begin accepting traffic
func (r *Runtime) Run() (err error) {
	wg := sync.WaitGroup{}
	wg.Add(len(r.daemons) + 2)

	go func() {
		defer wg.Done()
		r.handleRequest()
	}()

	for _, d := range r.daemons {
		log.Log("Daemon '%s' running.", d.Name())
		go func(daemon Daemon) {
			defer wg.Done()
			err = daemon.Run()
			if err != nil {
				log.Err(err)
				return
			}
		}(d)
	}

	delay := 1000.0
	for {
		select {
		case <-time.After(time.Millisecond * time.Duration(delay)):
			log.Event("Update timed out")
		default:
			pulse.Fixed(int(delay))
			start := time.Now()
			err = r.Update()
			if err != nil {
				log.ErrF(err, "runtime update error: %s")
			}
			d := time.Since(start)
			pulse.End()
			dur := (time.Millisecond * time.Duration(delay)) - d
			if dur > 0 {
				time.Sleep(dur)
			}

		}
		_ = r.Endpoints.Timings()

	}

}
