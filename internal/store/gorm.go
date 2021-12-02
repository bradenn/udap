// Copyright (c) 2021 Braden Nicholson

package store

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB
var client *redis.Client

type Database struct {
	database *gorm.DB
	client   *redis.Client
}

// Dependency is the level at which this service needs to run
func (d *Database) Dependency() (level int) {
	return 0
}

// Name is the name of the struct
func (d *Database) Name() (name string) {
	return "database"
}

// Load configures and prepares the parent struct for running
func (d *Database) Load() (err error) {
	pg := postgres.Open(dbURL())
	d.database, err = gorm.Open(pg, &gorm.Config{})
	if err != nil {
		return err
	}
	d.client = redis.NewClient(&redis.Options{
		Addr:      "localhost:6379",
		Password:  "", // no password set
		DB:        0,  // use default DB
		OnConnect: d.redisConnect,
	})
	fmt.Println(d.client.Get(context.Background(), "apple"))
	// Return no errors
	return nil
}

// Run will begin the main-sequence activities of the parent struct
func (d *Database) redisConnect(ctx context.Context, cn *redis.Conn) error {
	client = d.client
	return nil
}

func (d *Database) Run(interface{}) (err error) {
	DB = d.database
	return nil
}

// Cleanup will begin the main-sequence activities of the parent struct
func (d *Database) Cleanup() (err error) {
	d.database = nil
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
