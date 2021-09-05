package server

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "gorm.io/driver/postgres"
	"net/url"
	"os"
)

// NewDatabase initializes gormDB and connects to a postgresql database
func NewDatabase() (database *gorm.DB, err error) {
	// Using a preformatted URL, the database object is created
	database, err = gorm.Open("postgres", dbURL())
	if err != nil {
		return database, err
	}
	return database, err
}

// dbURL returns a formatted postgresql connection string.
func dbURL() string {
	// The credentials are retrieved from the OS environment
	dbUser := os.Getenv("dbUser")
	dbPass := os.Getenv("dbPass")
	// This credential pair is formed for later authentication
	user := url.UserPassword(dbUser, dbPass)
	// Host and port are also obtained from the environment
	dbHost := os.Getenv("dbHost")
	dbPort := os.Getenv("dbPort")
	// A connection string is built
	host := fmt.Sprintf("%s:%s", dbHost, dbPort)
	// The name of the database is again retrieved from the environment
	dbName := os.Getenv("dbName")
	// All variables are aggregated into the connection url
	u := url.URL{
		User:     user,
		Scheme:   "postgres",
		Host:     host,
		Path:     dbName,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	return u.String()
}
