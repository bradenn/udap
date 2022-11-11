// Copyright (c) 2022 Braden Nicholson

package io

import (
	"net/http"
	"time"
)

type Config interface {
	Set(key string, value string) error
	Get(key string) error
}

type Request interface {
	Push(r *http.Request) chan<- *http.Response
}

type Pilot struct {
	queue    chan *http.Request
	reqQueue chan chan *http.Response
}

type RequestBody struct {
	request  *http.Request
	response chan<- *http.Response
}

func (p *Pilot) listen() {
	for {
		select {
		case req := <-p.queue:
			err := p.run(req, <-p.reqQueue)
			if err != nil {
				break
			}
		}
	}
}

// Push adds a http request to the outgoing queue
func (p *Pilot) run(r *http.Request, re chan<- *http.Response) error {
	cli := http.Client{}
	cli.Timeout = time.Second * 2
	defer cli.CloseIdleConnections()

	do, err := cli.Do(r)
	if err != nil {
		return err
	}

	re <- do

	return nil
}

// Push adds a http request to the outgoing queue
func (p *Pilot) Push(r *http.Request) chan<- *http.Response {
	p.queue <- r
	return <-p.reqQueue
}

type IO struct {
	Services
	Request
	Config
}
