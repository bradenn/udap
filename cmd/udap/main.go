// Copyright (c) 2022 Braden Nicholson

package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"udap/internal/log"
	"udap/internal/orchestrator"
)

const VERSION = "2.17-beta"

func main() {

	err := setup()
	if err != nil {
		return
	}

	// Initialize Orchestrator
	o, err := orchestrator.NewOrchestrator()
	if err != nil {
		log.Err(err)
		return
	}

	// Initialize services
	err = o.Start()
	if err != nil {
		log.Err(err)
		return
	}

	// Run udap
	err = o.Run()
	if err != nil {
		log.Err(err)
		return
	}
}

func setup() error {
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
