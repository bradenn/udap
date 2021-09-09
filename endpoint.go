package main

import (
	"github.com/go-chi/chi"
	"net/http"
	"udap/server"
)

type Endpoint struct {
	Persistent
	Name    string `json:"name" gorm:"unique"`
	Enabled bool   `json:"enabled"`
}

func (e *Endpoint) Route(router chi.Router) {
	router.Post("/", createEndpoint)
	router.Get("/", findEndpoints)
	router.Get("/{id}", findEndpoint)
}

func createEndpoint(writer http.ResponseWriter, request *http.Request) {
	req, db := server.NewRequest(writer, request)

	var err error
	var model Endpoint

	req.DecodeModel(&model)
	model.Enabled = false

	err = db.Create(&model).Error
	if err != nil {
		req.Reject(err.Error(), http.StatusConflict)
		return
	}

	jwt, err := server.SignUUID(model.Id)
	if err != nil {
		req.Reject("Internal Error", http.StatusInternalServerError)
		return
	}

	resolve := map[string]interface{}{"token": jwt}

	req.Resolve(resolve, http.StatusOK)
}

func findEndpoints(writer http.ResponseWriter, request *http.Request) {
	req, db := server.NewRequest(writer, request)
	var model []Endpoint
	db.Model(&model).Find(&model)
	req.Resolve(model, http.StatusOK)
}

func findEndpoint(writer http.ResponseWriter, request *http.Request) {
	req, db := server.NewRequest(writer, request)

	var model Endpoint
	id := req.Param("id")
	db.Model(&model).Where("id = ?", id).Preload("Groups")

	err := db.Find(&model).Error
	if err != nil {
		req.Reject(err, http.StatusNotFound)
		return
	}

	req.Resolve(model, http.StatusOK)
}
