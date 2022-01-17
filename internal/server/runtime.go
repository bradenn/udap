// Copyright (c) 2021 Braden Nicholson

package server

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
	"udap/internal/bond"
	"udap/internal/controller"
	"udap/internal/log"
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
		start := time.Now()
		msg.Respond(r.ctrl.Handle(msg))
		log.Event("EVENT: %s.%s (%s)", msg.Target, msg.Operation, time.Since(start))
	}
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

	for _, d := range r.daemons {

		log.Log("Daemon '%s' loaded.", d.Name())
		err = d.Setup(r.ctrl, b)
		if err != nil {
			return
		}

	}

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
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {

	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func MacFromIpv4(ipv4 string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range interfaces {
		a, err := i.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range a {
			if addr.String() == fmt.Sprintf("%s/24", ipv4) {
				return i.HardwareAddr.String(), nil
			}
		}
	}
	return "", nil
}

type System struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Ipv4        string `json:"ipv4"`
	Ipv6        string `json:"ipv6"`
	Hostname    string `json:"hostname"`
	Mac         string `json:"mac"`
	Go          string `json:"go"`
	Cores       int    `json:"cores"`
}

var SystemInfo System

func (r *Runtime) Load() (err error) {

	ipv4 := GetOutboundIP().String()

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	fromIpv4, err := MacFromIpv4(ipv4)
	if err != nil {
		return err
	}

	r.System = System{
		Name:        "UDAP",
		Version:     os.Getenv("version"),
		Environment: os.Getenv("environment"),
		Hostname:    hostname,
		Ipv4:        ipv4,
		Mac:         fromIpv4,
		Go:          runtime.Version(),
		Cores:       runtime.NumCPU(),
	}
	SystemInfo = r.System
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

type DaemonData struct {
	Tick int
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

	go func() {
		defer wg.Done()
		delay := 5000.0
		for {
			start := time.Now()
			err = r.Update()
			if err != nil {
				log.ErrF(err, "runtime update error: %s")
			}
			d := time.Since(start)
			// log.Event("Tick: %.3d threads, %.2f%% load, %s", runtime.NumGoroutine(),
			// 	float64(d.Milliseconds())/delay, d.String())
			select {
			case <-time.After(time.Millisecond * time.Duration(delay)):
				log.ErrF(fmt.Errorf("timed out main update loop"), "%s")
				continue
			default:
				time.Sleep(time.Millisecond*time.Duration(delay) - d)
			}
		}
	}()

	wg.Wait()

	return nil
}
