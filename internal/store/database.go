// Copyright (c) 2021 Braden Nicholson

package store

import (
	"fmt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

type Database struct {
	*gorm.DB
}

func (d *Database) Dependency() (level int) {
	return 1
}

func (d *Database) Name() (name string) {
	return "database"
}

func (d *Database) Load() (err error) {
	pg := postgres.Open(dbURL())
	d.DB, err = gorm.Open(pg, &gorm.Config{})
	DB = d.DB
	if err != nil {
		return err
	}
	return nil
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
