package main

import (
	"context"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"udap/server"
)

type Endpoint struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Permissions []string           `json:"permissions" bson:"permissions"`
	Instances   []Instance         `json:"instances"`
	Enabled     bool               `json:"enabled"`
}

func (e *Endpoint) instances() error {
	var instances []Instance

	ctx := context.Background()

	find, err := server.Collection("instances").Find(ctx, bson.M{"permission": bson.M{"$in": e.Permissions}})
	if err != nil {
		return nil
	}

	for find.Next(ctx) {
		var instance Instance
		if err = find.Decode(&instance); err != nil {
			return nil
		}
		instances = append(instances, instance)
	}

	e.Instances = instances
	return nil
}

func (e *Endpoint) BeforeCreate() error {
	e.Enabled = false

	return nil
}

func (e *Endpoint) Route(router chi.Router) {
	router.Post("/", createEndpoint)
	router.Get("/", findEndpoints)
	router.Get("/{id}", findEndpoint)
}

func (e *Endpoint) AfterFind() error {

	err := e.instances()
	if err != nil {
		return err
	}
	return nil
}

func createEndpoint(writer http.ResponseWriter, request *http.Request) {
	req, ctx, db := server.NewRequest(writer, request, "endpoints")

	var err error
	var model Endpoint

	req.DecodeModel(&model)

	if err = model.BeforeCreate(); err != nil {
		return
	}

	result, err := db.InsertOne(ctx, &model)
	if err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	model.Id = result.InsertedID.(primitive.ObjectID)

	jwt, err := server.SignUUID(model.Id.Hex())
	if err != nil {
		req.Reject("Internal Error", http.StatusInternalServerError)
		return
	}

	resolve := map[string]interface{}{"token": jwt}

	req.Resolve(resolve, http.StatusOK)
}

func findEndpoints(writer http.ResponseWriter, request *http.Request) {
	req, ctx, db := server.NewRequest(writer, request, "endpoints")
	var model []Endpoint

	result, err := db.Find(ctx, bson.M{})
	if err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	for result.Next(ctx) {
		var endpoint Endpoint
		if err = result.Decode(&endpoint); err != nil {
			return
		}

		if err = endpoint.AfterFind(); err != nil {
			return
		}

		model = append(model, endpoint)
	}

	req.Resolve(model, http.StatusOK)
}

func findEndpoint(writer http.ResponseWriter, request *http.Request) {
	req, ctx, db := server.NewRequest(writer, request, "endpoints")

	id, err := req.ParamObjectId("id")
	if err != nil {
		return
	}

	var model Endpoint

	if err = db.FindOne(ctx, bson.M{"_id": id}).Decode(&model); err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	err = model.AfterFind()
	if err != nil {
		return
	}
	req.Resolve(model, http.StatusOK)
}
