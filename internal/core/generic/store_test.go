// Copyright (c) 2022 Braden Nicholson

package generic

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

var DB *gorm.DB
var mock sqlmock.Sqlmock

var store Store[Mock]

func TestNewStore(t *testing.T) {
	var db *sql.DB
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open mock sql db, got error: %v", err)
	}

	if mock == nil {
		t.Errorf("Mock is null: %v", err)
	}

	if db == nil {
		t.Errorf("Database is null: %v", err)
	}

	dial := postgres.New(postgres.Config{
		DSN:                  "mock",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	DB, err = gorm.Open(dial, &gorm.Config{})
	if err != nil {
		t.Errorf("Failed to open gorm v2 db, got error: %v", err)
	}

	store = NewStore[Mock](DB)

}

func TestFindById(t *testing.T) {

	elem, err := store.FindById("123")
	if err != nil {
		t.Errorf("Failed to FindById: %v", err)
	}

	fmt.Println(elem)

}
