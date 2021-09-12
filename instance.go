package main

import (
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
	"udap/server"
)

type Instance struct {
	Persistent
	Name        string    `json:"name" gorm:"unique"`
	Description string    `json:"description"`
	ModuleId    uuid.UUID `json:"moduleId"`
	Module      Module    `json:"module"`
}

func (i *Instance) Route(router chi.Router) {
	router.Post("/", createInstance)
	router.Get("/{id}", findInstance)
	router.Get("/{id}/func/{function}", runFunction)
	router.Get("/", findInstances)
}

func runFunction(w http.ResponseWriter, r *http.Request) {
	req, db := server.NewRequest(w, r)

	var model Instance

	id := req.Param("id")

	err := db.Where("id = ?", id).Preload("Module").First(&model).Error
	if err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	module, err := model.Module.Initialize(model.Id)
	if err != nil {
		return
	}

	function := req.Param("function")

	run, err := module.Run(function)
	if err != nil {
		return
	}

	req.ResolveRaw(run, http.StatusOK)
}

func findInstance(w http.ResponseWriter, r *http.Request) {
	req, db := server.NewRequest(w, r)

	var model Instance

	id := req.Param("id")

	err := db.Where("id = ?", id).Preload("Module").First(&model).Error
	if err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	req.Resolve(model, http.StatusOK)
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
	db.Model(&model).Preload("Module").Find(&model)
	req.Resolve(model, http.StatusOK)
}
