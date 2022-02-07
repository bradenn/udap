// Copyright (c) 2021 Braden Nicholson

package store

import (
	"fmt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

var DB Database

type Persistent struct {
	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"updated"`
	deletedAt *time.Time `sql:"index"`
	// Id is primary key of the persistent type, represented as a UUIDv4
	Id string `json:"id" gorm:"primary_key;type:string;default:uuid_generate_v4()"`
}

type Database struct {
	*gorm.DB
}

func NewDatabase() (Database, error) {
	pg := postgres.Open(dbURL())
	db, err := gorm.Open(pg, &gorm.Config{})
	if err != nil {
		return Database{}, err
	}

	DB.DB = db

	return DB, nil
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
