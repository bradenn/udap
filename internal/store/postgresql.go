// Copyright (c) 2021 Braden Nicholson

package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"reflect"
	"strings"
	"time"
)

var Pg *Postgresql

type Postgresql struct {
	*pgxpool.Pool
}

func (d *Postgresql) Dependency() (level int) {
	return 1
}

func (d *Postgresql) Name() (name string) {
	return "postgres"
}

type Wow struct {
	Persistent
}

type Object interface {
}

func (d *Postgresql) Insert() {

}

type Entity struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Module string `json:"module"`
	State  string `gorm:"-" json:"state"`
}

func (e *Entity) FromMap(data map[string]interface{}) {
	e.Id = data["id"].(string)
	e.Name = data["name"].(string)
	e.Type = data["type"].(string)
	e.Module = data["module"].(string)
	e.State = data["state"].(string)
}

func (e *Entity) JSON() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func (d *Postgresql) Load() (err error) {
	d.Pool, err = pgxpool.Connect(context.Background(), pqUrl())
	Pg = d
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	t := time.Now()
	eT := reflect.TypeOf(Entity{})
	var fields []string
	for _, field := range reflect.VisibleFields(eT) {
		fields = append(fields, field.Tag.Get("json"))
	}
	query, err := d.Query(context.Background(), fmt.Sprintf("select %s from entities where true", strings.Join(fields,
		", ")))
	if err != nil {
		return err
	}

	for query.Next() {

		values, err := query.Values()
		if err != nil {
			return err
		}
		var res map[string]interface{}
		res = map[string]interface{}{}
		for i, field := range fields {
			res[field] = values[i]
		}
		entity := Entity{}
		entity.FromMap(res)
		fmt.Println(entity.JSON())
	}

	fmt.Println(time.Since(t))
	return nil
}

// dbURL returns a formatted postgresql connection string.
func pqUrl() string {
	// The credentials are retrieved from the OS environment
	dbUser := os.Getenv("dbUser")
	dbPass := os.Getenv("dbPass")
	// Host and port are also obtained from the environment
	dbHost := os.Getenv("dbHost")
	dbPort := os.Getenv("dbPort")
	// The name of the database is again retrieved from the environment
	dbName := os.Getenv("dbName")
	// All variables are aggregated into the connection url
	u := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	return u
}
