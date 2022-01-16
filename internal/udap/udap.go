// Copyright (c) 2021 Braden Nicholson

package udap

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"udap/internal/cache"
	"udap/internal/log"
	"udap/internal/models"
	"udap/internal/server"
	"udap/internal/store"
)

const VERSION = "2.8.2"

type Udap struct {
	runtime  *server.Runtime
	cache    cache.Cache
	database store.Database
}

func Run() error {
	err := config()
	if err != nil {
		return nil
	}

	u := &Udap{}

	u.cache, err = cache.NewCache()
	if err != nil {
		return err
	}

	u.database, err = store.NewDatabase()
	if err != nil {
		return err
	}

	err = u.migrate()
	if err != nil {
		return err
	}

	u.runtime = &server.Runtime{}

	err = u.runtime.Load()
	if err != nil {
		return err
	}

	err = u.runtime.Run()
	if err != nil {
		return err
	}
	return nil
}

func (u *Udap) migrate() error {
	err := u.database.AutoMigrate(models.Log{}, models.Endpoint{}, models.Entity{}, models.Module{}, models.Device{},
		models.Network{})
	if err != nil {
		return err
	}
	return nil
}

func config() error {

	log.Log("UDAP v%s - Copyright (c) 2021 Braden Nicholson", VERSION)
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load .env file could not find any environment variables")
	}
	if os.Getenv("environment") == "production" {
		log.Log("Running in PRODUCTION mode.")
	} else {
		log.Log("Running in DEVELOPMENT mode.")
	}
	err = os.Setenv("version", VERSION)
	if err != nil {
		return err
	}
	return nil
}

// // Go is run to begin the program
// func (u *Udap) Go() (err error) {
// 	var dependencies []Dependency
// 	// We pull the services from the map
// 	for _, dependency := range u.dependencies {
// 		dependencies = append(dependencies, *dependency)
// 	}
// 	// Next we sort the services by their priorities
// 	sort.Slice(dependencies, func(i, j int) bool {
// 		return dependencies[i].Dependency() < dependencies[j].Dependency()
// 	})
//
// 	for _, dep := range dependencies {
// 		t := time.Now()
// 		err = dep.Load()
// 		if err != nil {
// 			log.Err(err)
// 		}
// 		log.Log("Dependency '%s' loaded. (%s)", dep.Name(), time.Since(t).String())
//
// 	}
//
// 	log.Log("All dependencies loaded")
// 	// A wait group is made to prevent premature exit
// 	wg := sync.WaitGroup{}
// 	// Each service is given a slot in the wait group
// 	wg.Add(len(u.services))
// 	// This array will contain a list of services
// 	var services []Service
// 	err = u.migrate()
// 	if err != nil {
// 		return err
// 	}
// 	// We pull the services from the map
// 	for _, service := range u.services {
// 		services = append(services, *service)
// 	}
// 	// Next we sort the services by their priorities
// 	sort.Slice(services, func(i, j int) bool {
// 		return services[i].Dependency() < services[j].Dependency()
// 	})
// 	sg := sync.WaitGroup{}
// 	sg.Add(len(services))
// 	// Now each service runs, with the priorities 1 and 0 blocking
// 	for _, s := range services {
// 		go func(service Service) {
// 			log.Log("Service '%s' running", service.Name())
// 			sg.Done()
// 			err = service.Run()
// 			if err != nil {
// 				log.Err(err)
// 				return
// 			}
// 			wg.Done()
// 		}(s)
//
// 	}
//
// 	sg.Wait()
// 	log.Log("Running.")
//
// 	wg.Wait()
// 	return nil
// }
