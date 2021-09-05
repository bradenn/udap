package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func Init() {
	if os.Getenv("ENV") == "production" {
		fmt.Println("Running in PRODUCTION mode.")
	} else {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Failed to load .env file could not find any environment variables.")
		}
		fmt.Println("Running in DEVELOPMENT mode.")
	}
}
