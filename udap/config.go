// Copyright (c) 2021 Braden Nicholson

package udap

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func Init() {
	if os.Getenv("ENV") == "production" {
		Info("Running in PRODUCTION mode.")
	} else {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Failed to load .env file could not find any environment variables.")
		}
		Info("Running in DEVELOPMENT mode.")
	}
}
