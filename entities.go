package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
	"udap/server"
)

type Instance struct {
	Persistent
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
	Permission  string `json:"permission"`
	Module      string `json:"module"`

	environment string // Private, unless polled manually by function
}

func (i *Instance) env() (o Object) {
	err := json.Unmarshal([]byte(i.environment), &o)
	if err != nil {
		return nil
	}
	return o
}

func (i *Instance) Route(router chi.Router) {
	router.Post("/", createInstance)
	router.Get("/", findInstances)
}

func createInstance(w http.ResponseWriter, r *http.Request) {
	req, db := server.NewRequest(w, r)
	var model Instance

	req.DecodeModel(&model)

	err := db.Create(&model).Error
	if err != nil {
		req.Reject(err.Error(), http.StatusConflict)
		return
	}

	req.Resolve(model, http.StatusOK)
}

func findInstances(w http.ResponseWriter, r *http.Request) {
	req, db := server.NewRequest(w, r)
	var model []Instance
	db.Model(&model).Find(&model)
	req.Resolve(model, http.StatusOK)
}
