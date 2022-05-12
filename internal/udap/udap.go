// Copyright (c) 2021 Braden Nicholson

package udap

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"udap/internal/log"
	"udap/internal/models"
	"udap/internal/server"
	"udap/internal/store"
)

const VERSION = "2.12"

type Udap struct {
	runtime *server.Runtime
}

func (u Udap) startup() error {
	log.Log("UDAP v%s - Copyright (c) 2019-2022 Braden Nicholson", VERSION)

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load .env file")
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

func Start() error {
	u := &Udap{}
	err := u.startup()
	if err != nil {
		return err
	}

	_, err = store.NewDatabase()
	if err != nil {
		return err
	}

	err = models.MigrateModels()
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
