package main

import (
	"net/http"
	"time"
	"udap/server"
	"udap/server/auth"
)

type EndpointResolve struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}

type Endpoint struct {
	Persistent
	Name    string `json:"name" gorm:"unique"`
	Enabled bool   `json:"enabled"`
}

func (e *Endpoint) Create(writer http.ResponseWriter, request *http.Request) {
	req, db := server.NewRequest(writer, request)

	var err error
	var model Endpoint

	err = req.DecodeModel(&model)
	if err != nil {
		req.Reject(err.Error(), http.StatusBadRequest)
		return
	}

	model.Enabled = false

	err = db.Create(&model).Error
	if err != nil {
		req.Reject(err.Error(), http.StatusConflict)
		return
	}

	tokenBody := map[string]interface{}{}
	tokenBody["id"] = model.Id.String()

	jwt, err := auth.SignKey(tokenBody)
	if err != nil {
		req.Reject("Internal Error", http.StatusInternalServerError)
		return
	}
	resolve := EndpointResolve{
		Token:     jwt,
		CreatedAt: time.Now(),
	}
	req.JSON(resolve, http.StatusOK)
}

func (e *Endpoint) FindAll(writer http.ResponseWriter, request *http.Request) {
	req, db := server.NewRequest(writer, request)

	var model []Endpoint

	db.Model(&model).Find(&model)

	req.JSON(model, http.StatusOK)
}

func (e *Endpoint) FindOne(writer http.ResponseWriter, request *http.Request) {
	req, db := server.NewRequest(writer, request)

	var model Endpoint

	id := req.Param("id")

	db.Model(&model).Where("id = ?", id).Preload("Groups")

	err := db.Find(&model).Error
	if err != nil {
		req.JSON(err, http.StatusNotFound)
		return
	}

	req.JSON(model, http.StatusOK)
}
