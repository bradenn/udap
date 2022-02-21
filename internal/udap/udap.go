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

const VERSION = "2.9.3"

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
	log.Log("UDAP v%s - Copyright (c) 2019-2022 Braden Nicholson", VERSION)
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
