// Copyright (c) 2021 Braden Nicholson

package udap

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"udap/internal/log"
	"udap/internal/models"
	"udap/internal/store"
)

func config() error {
	log.Log("UDAP v%s - Copyright (c) 2021 Braden Nicholson", VERSION)
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load .env file could not find any environment variables")
	}
	if os.Getenv("ENV") == "production" {
		log.Log("Running in PRODUCTION mode.")
	} else {
		log.Log("Running in DEVELOPMENT mode.")
	}
	return nil
}

func migrate() (err error) {
	err = store.DB.AutoMigrate(models.Endpoint{}, models.Entity{}, models.Module{},
		models.Instance{}, models.Subscription{}, models.Grant{})
	if err != nil {
		return err
	}
	return nil
}
