package server

import (
	"fmt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func NewGormDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbURL()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, err
}

// dbURL returns a formatted postgresql connection string.
func dbURL() string {
	// The credentials are retrieved from the OS environment
	dbUser := os.Getenv("dbUser")
	dbPass := os.Getenv("dbPass")
	// Host and port are also obtained from the environment
	dbHost := os.Getenv("dbHost")
	dbPort := os.Getenv("dbPort")
	// The name of the database is again retrieved from the environment
	dbName := os.Getenv("dbName")
	// All variables are aggregated into the connection url
	u := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbHost, dbUser, dbPass, dbName, dbPort)
	return u
}
