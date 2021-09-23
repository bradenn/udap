package main

import (
	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"net/http"
	"udap/server"
)

type Instance struct {
	Persistent
	Permission  string   `json:"permission" bson:"permission"`
	Functions   []string `json:"functions" bson:"-"`
	Name        string   `json:"name" gorm:"unique"`
	Description string   `json:"description"`
}

func (i *Instance) AfterFind(_ *gorm.DB) error {
	return nil
}

func (i *Instance) BeforeCreate(_ *gorm.DB) error {
	return nil
}

func RouteInstances(router chi.Router) {
	router.Post("/", createInstance)
	router.Get("/{id}", findInstance)
	router.Get("/{id}/func/{function}", runFunction)
	router.Get("/", findInstances)
}

func runFunction(w http.ResponseWriter, r *http.Request) {
	req, db := server.NewRequest(w, r)

	id := req.Param("id")

	db.Where("id = ?", id)
	if err := db.Error; err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	req.ResolveRaw("", http.StatusOK)
}

func findInstance(w http.ResponseWriter, r *http.Request) {
	req, db := server.NewRequest(w, r)

	id := req.Param("id")

	var model Instance

	db.Where("id = ?", id).First(&model)

	if err := db.Error; err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	req.Resolve(model, http.StatusOK)
}

func createInstance(w http.ResponseWriter, r *http.Request) {
	req, db := server.NewRequest(w, r)

	var model Instance

	req.DecodeModel(&model)
	db.Create(&model)

	if err := db.Error; err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	req.Resolve(model, http.StatusOK)
}

func findInstances(w http.ResponseWriter, r *http.Request) {
	req, db := server.NewRequest(w, r)
	var models []Endpoint

	db.Find(&models)

	if err := db.Error; err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	req.Resolve(models, http.StatusOK)
}
