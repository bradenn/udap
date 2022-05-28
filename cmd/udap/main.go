// Copyright (c) 2022 Braden Nicholson

package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"udap/internal/log"
	"udap/internal/orchestrator"
)

const VERSION = "2.13"

func main() {

	err := setup()
	if err != nil {
		return
	}

	// Initialize Orchestrator
	o := orchestrator.NewOrchestrator()

	// Initialize services
	err = o.Init()
	if err != nil {
		return
	}

	// Run udap
	err = o.Run()
	if err != nil {
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
