// Copyright (c) 2022 Braden Nicholson

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
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

func (w *Worldspace) Run() error {
	w.server = http.Server{}
	w.server.Addr = "0.0.0.0:5058"
	router := chi.NewRouter()
	err := w.Triggers.Register(&domain.Trigger{
		Name:        "motion",
		Type:        "module",
		Description: "When motion is detected",
	})
	if err != nil {
		return err
	}
	router.Post("/motion/{zone}", w.handleMotion)
	err = w.Triggers.Register(&domain.Trigger{
		Name:        "homekit-arrival",
		Type:        "module",
		Description: "When someone arrives",
	})
	if err != nil {
		return err
	}
	router.Post("/arrive/{id}", w.handleArrival)
	err = w.Triggers.Register(&domain.Trigger{
		Name:        "homekit-departure",
		Type:        "module",
		Description: "When someone departs",
	})
	if err != nil {
		return err
	}
	router.Post("/depart/{id}", w.handleDeparture)
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
