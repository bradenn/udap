package main

import (
	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"net/http"
	"udap/server"
)

type Endpoint struct {
	Persistent
	Name        string     `json:"name" bson:"name"`
	Permissions []string   `json:"permissions" bson:"permissions"`
	Instances   []Instance `json:"instances"`
	Enabled     bool       `json:"enabled"`
}

func (e *Endpoint) BeforeCreate(_ *gorm.DB) error {
	e.Enabled = false

	return nil
}

func RouteEndpoints(router chi.Router) {
	router.Post("/", createEndpoint)
	router.Get("/", findEndpoints)
	router.Get("/{id}", findEndpoint)
	router.Get("/poll", pollEndpoint)
}

func (e *Endpoint) AfterFind(_ *gorm.DB) error {
	return nil
}

func createEndpoint(writer http.ResponseWriter, request *http.Request) {
	req, db := server.NewRequest(writer, request)

	var err error
	var model Endpoint

	req.DecodeModel(&model)
	db.Create(&model)

	if err := db.Error; err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	jwt, err := server.SignUUID(model.Id.String())
	if err != nil {
		req.Reject("Internal Error", http.StatusInternalServerError)
		return
	}

	resolve := map[string]interface{}{"token": jwt}

	req.Resolve(resolve, http.StatusOK)
}

func findEndpoints(writer http.ResponseWriter, request *http.Request) {
	req, db := server.NewRequest(writer, request)
	var models []Endpoint

	db.Find(&models)

	if err := db.Error; err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	req.Resolve(models, http.StatusOK)
}

func findEndpoint(writer http.ResponseWriter, request *http.Request) {
	req, db := server.NewRequest(writer, request)

	id := req.Param("id")

	var model Endpoint

	db.Where("id = ?", id).First(&model)

	if err := db.Error; err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	req.Resolve(model, http.StatusOK)
}

func pollEndpoint(writer http.ResponseWriter, request *http.Request) {
	req, db := server.NewRequest(writer, request)

	id := req.JWTClaim("id").(string)

	var model Endpoint

	db.Where("id = ?", id).First(&model)

	if err := db.Error; err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	req.Resolve(model, http.StatusOK)
}
