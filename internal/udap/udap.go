// Copyright (c) 2021 Braden Nicholson

package udap

import (
	"context"
	"sort"
	"sync"
	"udap/internal/log"
	"udap/internal/server"
)

const VERSION = "2.4.5"

var udap *Udap

type Udap struct {
	services map[string]Service
	runtime  *server.Runtime
}

func New() *Udap {
	err := config()
	if err != nil {
		return nil
	}
	return &Udap{
		services: map[string]Service{},
	}
}

func (u *Udap) Add(service ...Service) {
	for _, s := range service {
		u.services[s.Name()] = s
	}
}

// Go is run to begin the program
func (u *Udap) Go() (err error) {

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
	var names []string
	// Now each service runs, with the priorities 1 and 0 blocking
	for _, service := range services {
		names = append(names, service.Name())
		level := service.Dependency()
		if level == 0 || level == 1 {
			goService(service, &wg)
			log.Sherlock("Service '%s' loaded", service.Name())
		} else {
			go goService(service, &wg)
			log.Sherlock("Service '%s' running", service.Name())
		}

	}
	err = migrate()
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func (u *Udap) load() (err error) {

	return nil
}

func goService(service Service, wg *sync.WaitGroup) {
	var err error
	defer func() {
		err = service.Cleanup()
		if err != nil {
			log.Err(err)
		}
		wg.Done()
	}()
	err = service.Load()
	if err != nil {
		log.Err(err)
	}
	ctx := context.Background()

	err = service.Run(ctx)
	if err != nil {
		log.Err(err)
	}
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

func (u *Udap) cleanup() (err error) {
	for _, service := range u.services {
		err = service.Cleanup()
		if err != nil {
			return err
		}
	}
	return nil
}

type Service interface {
	Name() (name string)
	Load() (err error)
	Dependency() (level int)
	Run(ctx interface{}) (err error)
	Cleanup() (err error)
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
