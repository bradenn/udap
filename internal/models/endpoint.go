// Copyright (c) 2021 Braden Nicholson

package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"math/rand"
	"time"
	"udap/internal/log"
	"udap/internal/store"
)

type Connection struct {
	WS     *websocket.Conn
	active *bool
	edit   chan any
	done   chan bool
}

func (c *Connection) Active() bool {
	return *c.active
}

func (c *Connection) Send(body any) {
	if c.Active() {
		c.edit <- body
	}
}

func NewConnection(ws *websocket.Conn) *Connection {
	ch := make(chan any)
	d := make(chan bool)
	a := true
	c := &Connection{
		WS:     ws,
		edit:   ch,
		done:   d,
		active: &a,
	}

	return c
}

func (c *Connection) Close() {
	close(c.edit)
	a := false
	c.active = &a
}

func (c *Connection) Watch() {
	for a := range c.edit {
		if c.WS == nil {
			return
		}
		err := c.WS.WriteJSON(a)
		if err != nil {
			continue
		}
	}
}

// Endpoint represents a client device connected to the UDAP network
type Endpoint struct {
	store.Persistent

	Name string `json:"name" gorm:"unique"`

	Type string `json:"type" gorm:"default:'terminal'"`

	Frequency int `json:"frequency" gorm:"default:3000"`

	Connected bool `json:"connected"`

	Key string `json:"key"`

	registered    bool
	Connection    *Connection `json:"-" gorm:"-"`
	enrolledSince time.Time   `gorm:"-"`
}

func (e *Endpoint) Compile() (a map[string]any) {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(marshal, &a)
	if err != nil {
		return nil
	}
	return a
}

func (e *Endpoint) Enroll(ws *websocket.Conn) error {
	err := store.DB.Model(&Endpoint{}).FirstOrCreate(&e).Error
	if err != nil {
		return err
	}
	ws.SetCloseHandler(e.closeHandler)
	e.Connection = NewConnection(ws)
	e.Connected = true
	e.registered = true
	e.enrolledSince = time.Now()
	log.Log("Endpoint '%s' enrolled (%s)", e.Name, ws.LocalAddr())

	return nil
}

func NewEndpoint(name string) Endpoint {
	endpoint := Endpoint{}
	endpoint.Name = name
	endpoint.Type = "terminal"
	return endpoint
}

func (e *Endpoint) closeHandler(code int, text string) error {
	if e.Enrolled() {
		e.Unenroll()
	}

	return nil
}

func (e *Endpoint) Enrolled() bool {
	return e.registered
}

// BeforeCreate is a hook function from gorm, called when an endpoint is inserted
func (e *Endpoint) BeforeCreate(_ *gorm.DB) error {
	e.Key = randomSequence()
	return nil
}

// randomSequence generates a random id for use as a key
func randomSequence() string {
	template := "abcdefghijklmnopqrstuvwxyz"
	var out string
	rand.Seed(time.Now().Unix())
	for i := 0; i < 8; i++ {
		r := rand.Intn(26)
		u := template[r]
		out += string(u)
	}
	return out
}

// Fetch populates the struct from the database
func (e *Endpoint) Fetch() error {
	if e.Id == "" {
		return fmt.Errorf("invalid id")
	}
	err := store.DB.Where("name = ? OR id = ?", e.Name, e.Id).First(&e).Error
	if err != nil {
		return err
	}

	return nil
}

func (e *Endpoint) Unenroll() {
	e.registered = false
	e.Connected = false
	e.Connection.Close()
	log.Log("Endpoint '%s' unenrolled (%s)", e.Name, time.Since(e.enrolledSince).String())
}
