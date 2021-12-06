// Copyright (c) 2021 Braden Nicholson

package udap

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
	"udap/internal/cache"
	"udap/internal/log"
	"udap/internal/server"
	"udap/internal/store"
)

const VERSION = "2.4.5"

type Udap struct {
	services     map[string]Service
	dependencies map[string]Dependency
	ctx          context.Context
	cache.Cache
	store.Database
	server.Server
	server.Runtime
}

func Run() (*Udap, error) {
	err := config()

	if err != nil {
		return nil, err
	}

	u := &Udap{
		services:     map[string]Service{},
		dependencies: map[string]Dependency{},
		ctx:          context.Background(),
	}

	u.AddDependency(&u.Database, &u.Server, &u.Cache, &u.Runtime)
	u.AddService(&u.Runtime, &u.Server)

	u.ctx = context.WithValue(u.ctx, "runtime", &u.Runtime)
	err = u.Go()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *Udap) AddService(service ...Service) {
	for _, s := range service {
		u.services[s.Name()] = s
	}
}

func (u *Udap) AddDependency(dependencies ...Dependency) {
	for _, s := range dependencies {
		u.dependencies[s.Name()] = s
	}
}

// Go is run to begin the program
func (u *Udap) Go() (err error) {

	dg := sync.WaitGroup{}
	dg.Add(len(u.dependencies))
	log.Proc("Loading dependencies: ")

	for _, dep := range u.dependencies {
		fmt.Printf("%s ", dep.Name())

		go func(dependency Dependency) {
			t := time.Now()
			err := dependency.Load()
			if err != nil {
				log.Err(err)
			}
			u.ctx = context.WithValue(u.ctx, dependency.Name(), dependency)

			log.Log("Dependency '%s' loaded. (%s)", dependency.Name(), time.Since(t).String())
			dg.Done()
		}(dep)
	}
	fmt.Println()
	dg.Wait()
	log.Log("All dependencies loaded")
	// A wait group is made to prevent premature exit
	wg := sync.WaitGroup{}
	// Each service is given a slot in the wait group
	wg.Add(len(u.services))
	// This array will contain a list of services
	var services []Service
	// We pull the services from the map
	for _, service := range u.services {
		services = append(services, service)
	}
	// Next we sort the services by their priorities
	sort.Slice(services, func(i, j int) bool {
		return services[i].Dependency() < services[j].Dependency()
	})
	sg := sync.WaitGroup{}
	sg.Add(len(services))
	// Now each service runs, with the priorities 1 and 0 blocking
	for _, s := range services {
		go func(service Service) {
			log.Log("Service '%s' running", service.Name())
			sg.Done()
			defer func() {
				// We panicked. Don't Panic!
				if r := recover(); r != nil {
					log.Log("Panic Recovered: %x", r)
				}
			}()
			err = service.Run(u.ctx)
			if err != nil {
				log.Err(err)
				return
			}
			wg.Done()
		}(s)

	}
	err = u.migrate()
	if err != nil {
		return err
	}

	sg.Wait()
	log.Log("Running.")
	wg.Wait()
	return nil
}

func (u *Udap) load() (err error) {

	return nil
}

func (u *Udap) run() (err error) {
	log.Sherlock("Now running.")
	wg := sync.WaitGroup{}
	wg.Add(len(u.services))
	for _, service := range u.services {
		s := service
		go func() {
			err = s.Run(context.Background())
			if err != nil {
				log.Err(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

type Core interface {
	Name() (name string)
	Dependency() (level int)
}

type Dependency interface {
	Core
	Load() (err error)
}

type Service interface {
	Core
	Run(ctx context.Context) (err error)
}

type DefaultService struct {
}

func (d DefaultService) Name() (name string) {
	return ""
}

func (d DefaultService) Load() (err error) {
	return nil
}

func (d DefaultService) Run() (err error) {
	return nil
}

func (d DefaultService) Cleanup() (err error) {
	return nil
}
