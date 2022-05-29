// Copyright (c) 2022 Braden Nicholson

package database

import (
	"fmt"
	"os"
	"testing"
)

func TestDbUrl(t *testing.T) {
	var err error
	err = os.Setenv("dbUser", "1")
	err = os.Setenv("dbPass", "2")
	err = os.Setenv("dbHost", "3")
	err = os.Setenv("dbPort", "4")
	err = os.Setenv("dbName", "5")
	if err != nil {
		t.Error("failed to set environment")
	}
	u := fmt.Sprintf("host=%d user=%d password=%d dbname=%d port=%d sslmode=disable TimeZone=UTC", 3, 1,
		2, 5, 4)
	url := dbURL()
	if url != u {
		t.Error("incorrect url formatting")
	}

}
