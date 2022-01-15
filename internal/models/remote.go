// Copyright (c) 2022 Braden Nicholson

package models

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
	"udap/internal/log"
)

type Remote struct {
	conn *websocket.Conn
	rw   *sync.Mutex
}

func NewRemote(c *websocket.Conn) Remote {
	r := Remote{
		rw: &sync.Mutex{},
	}
	r.rw.Lock()
	r.conn = c
	log.Log("Get socket opened: %s", c.RemoteAddr())
	r.conn.SetCloseHandler(r.closeHandler)
	r.rw.Unlock()
	return r
}

func (r *Remote) closeHandler(code int, text string) error {
	if text == "" {
		text = "[empty]"
	}
	log.Log("Get socket closed: %s (%d)", text, code)
	return nil
}

func (r *Remote) Send(body json.RawMessage) error {
	r.rw.Lock()
	err := r.conn.WriteJSON(body)
	if err != nil {
		log.Err(err)
	}
	r.rw.Unlock()
	return err
}

func (r *Remote) Close() error {

	err := r.conn.Close()
	if err != nil {
		log.Err(err)
	}
	return err
}

func (r *Remote) Read() (body json.RawMessage, err error) {
	r.rw.Lock()
	err = r.conn.ReadJSON(body)
	if err != nil {
		log.Err(err)
	}
	r.rw.Unlock()
	return body, err
}
