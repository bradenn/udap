// Copyright (c) 2021 Braden Nicholson

package sql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
)

var conn *pgx.Conn

func init() {
	var err error
	conn, err = pgx.Connect(context.Background(), "postgres://udap:udap@127.0.0.1:5432/udap")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

}

type DBObject interface {
	TableName() string
}

func Select(destination DBObject) {

}

func Insert() {

}

func Update() {

}

func Delete() {

}
