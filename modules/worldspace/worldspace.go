// Copyright (c) 2022 Braden Nicholson

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
)

var Module Worldspace

type Worldspace struct {
	plugin.Module
	server   http.Server
	entityId string
}

type Landmarks struct {
	RightEye struct {
		Xa int `json:"xa"`
		Ya int `json:"ya"`
		Xb int `json:"xb"`
		Yb int `json:"yb"`
	} `json:"rightEye"`
	LeftEye struct {
		Xa int `json:"xa"`
		Ya int `json:"ya"`
		Xb int `json:"xb"`
		Yb int `json:"yb"`
	} `json:"leftEye"`
	Nose struct {
		Xa int `json:"x"`
		Ya int `json:"y"`
	} `json:"nose"`
}

type Prediction struct {
	Name      string    `json:"name"`
	Top       int       `json:"top"`
	Right     int       `json:"right"`
	Bottom    int       `json:"bottom"`
	Left      int       `json:"left"`
	Distance  float64   `json:"distance"`
	Landmarks Landmarks `json:"landmarks"`
}

type Status struct {
	Zone        string       `json:"zone"`
	Predictions []Prediction `json:"predictions"`
	Updated     time.Time    `json:"updated"`
}

func init() {
	config := plugin.Config{
		Name:        "worldspace",
		Type:        "module",
		Description: "worldspace integration",
		Version:     "0.1.1",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
}

func (w *Worldspace) Setup() (plugin.Config, error) {
	err := w.UpdateInterval(2000)
	if err != nil {
		return plugin.Config{}, err
	}
	return w.Config, nil
}

func (w *Worldspace) Update() error {
	if w.Ready() {

	}
	return nil
}

type Departure struct {
	Meta string `json:"meta"`
	Time string `json:"time"`
}

func (w *Worldspace) handleDeparture(writer http.ResponseWriter, request *http.Request) {
	// id := chi.URLParam(request, "id")
	a := Departure{}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf.Bytes(), &a)
	if err != nil {
		return
	}

	// err := w.Attributes.Update(w.entityId, "homeKitArrivals", string(marshal), time.Now())
	// if err != nil {
	// 	w.LogF("%s", err.Error())
	// 	return
	// }

	err = w.Triggers.Trigger("homekit-departure")
	if err != nil {
		return
	}
	writer.WriteHeader(200)

}

type Arrival struct {
	Meta string `json:"meta"`
	Time string `json:"time"`
}

func (w *Worldspace) handleArrival(writer http.ResponseWriter, request *http.Request) {
	// id := chi.URLParam(request, "id")
	a := Arrival{}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(buf.Bytes(), &a)
	if err != nil {
		return
	}

	// err := w.Attributes.Update(w.entityId, "homeKitArrivals", string(marshal), time.Now())
	// if err != nil {
	// 	w.LogF("%s", err.Error())
	// 	return
	// }
	err = w.Triggers.Trigger("homekit-arrival")
	if err != nil {
		return
	}

	writer.WriteHeader(200)

}

func (w *Worldspace) handleMotion(writer http.ResponseWriter, request *http.Request) {
	// zone := chi.URLParam(request, "zone")

	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		return
	}
	defer request.Body.Close()

	err = w.Triggers.Trigger("motion")
	if err != nil {
		return
	}

	// var p []Prediction
	//
	// err = json.Unmarshal(buf.Bytes(), &p)
	// if err != nil {
	// 	w.LogF("%s", err.Error())
	// 	return
	// }
	//
	// s := Status{
	// 	Zone:        zone,
	// 	Predictions: p,
	// 	Updated:     time.Now(),
	// }
	//
	// marshal, err := json.Marshal(s)
	// if err != nil {
	// 	w.LogF("%s", err.Error())
	// 	return
	// }
	//
	// err = w.Attributes.Update(w.entityId, "deskFace", string(marshal), time.Now())
	// if err != nil {
	// 	w.LogF("%s", err.Error())
	// 	return
	// }

	writer.WriteHeader(200)

}

func (w *Worldspace) endpointTrigger(router chi.Router, name string, desc string) error {
	router.Post(fmt.Sprintf("/%s", name), func(writer http.ResponseWriter, request *http.Request) {
		// zone := chi.URLParam(request, "zone")
		var buf bytes.Buffer
		_, err := buf.ReadFrom(request.Body)
		if err != nil {
			return
		}
		defer request.Body.Close()

		err = w.Triggers.Trigger(fmt.Sprintf("ws-%s", name))
		if err != nil {
			return
		}
		_, err = writer.Write([]byte("OK"))
		if err != nil {
			return
		}
	})

	err := w.Triggers.Register(&domain.Trigger{
		Name:        fmt.Sprintf("ws-%s", name),
		Type:        "module",
		Description: desc,
	})
	if err != nil {
		return err
	}
	return nil
}

func (w *Worldspace) Run() error {
	w.server = http.Server{}
	w.server.Addr = "0.0.0.0:5058"
	router := chi.NewRouter()
	err := w.endpointTrigger(router, "motion", "dummy motion")
	if err != nil {
		return err
	}

	err = w.endpointTrigger(router, "depart", "homekit depart")
	if err != nil {
		return err
	}

	err = w.endpointTrigger(router, "arrive", "homekit arrive")
	if err != nil {
		return err
	}

	err = w.endpointTrigger(router, "motion-1", "Living Room Motion 1")
	if err != nil {
		return err
	}

	err = w.endpointTrigger(router, "motion-2", "Living Room Motion 2")
	if err != nil {
		return err
	}

	w.server.Handler = router
	entity := domain.Entity{
		Name:   "faces",
		Type:   "media",
		Module: "worldspace",
	}

	err = w.Entities.Register(&entity)
	if err != nil {
		return err
	}

	w.entityId = entity.Id

	attr := domain.Attribute{
		Key:     "deskFace",
		Value:   "{}",
		Request: "{}",
		Type:    "media",
		Order:   0,
		Entity:  w.entityId,
	}

	err = w.Attributes.Register(&attr)
	if err != nil {
		return err
	}

	wsE := domain.Entity{
		Name:   "worldspace-1",
		Type:   "media",
		Module: "worldspace",
	}

	err = w.Entities.Register(&wsE)
	if err != nil {
		return err
	}
	path := make(chan domain.Attribute)
	dim := &domain.Attribute{
		Key:     "dim",
		Value:   "0",
		Request: "0",
		Type:    "range",
		Order:   0,
		Entity:  wsE.Id,
		Channel: path,
	}

	go func() {
		for attribute := range path {
			parseInt, err := strconv.Atoi(attribute.Request)
			if err != nil {
				w.Err(err)
				continue
			}

			cli := http.Client{
				Timeout: 500 * time.Millisecond,
			}
			post, err := cli.Post("http://10.0.1.85/dim", "application/json", bytes.NewReader([]byte(fmt.Sprintf("{ \"dim\": %d }", parseInt))))
			if err != nil {
				w.Err(err)
				continue
			}
			post.Body.Close()
			err = w.Attributes.Set(wsE.Id, "dim", attribute.Request)
			if err != nil {
				w.Err(err)
				continue
			}
		}
	}()

	err = w.Attributes.Register(dim)
	if err != nil {
		return err
	}
	go func() {
		err = w.server.ListenAndServe()
		if err != nil {
			return
		}
	}()
	return nil
}

func (w *Worldspace) Dispose() error {
	ctx := context.Background()
	err := w.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
