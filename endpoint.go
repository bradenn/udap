package main

import (
	"net/http"
	"udap/server"
)

//type EndpointToken struct {
//	Id        uuid.UUID `json:"id"`
//	CreatedAt time.Time `json:"createdAt"`
//}

type Endpoint struct {
	Persistent
	Name    string `json:"name" gorm:"unique"`
	Enabled bool   `json:"enabled"`
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
