package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"udap/server"
)

type Instance struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty" `
	Module      primitive.ObjectID `json:"module" bson:"module"`
	Permission  string             `json:"permission" bson:"permission"`
	Functions   []string           `json:"functions" bson:"-"`
	Name        string             `json:"name" gorm:"unique"`
	Description string             `json:"description"`
}

func (i *Instance) db() *mongo.Collection {
	return server.From("instances")
}

func (i *Instance) GetModule() (module Module, err error) {

	result := server.From("modules").FindOne(context.Background(), bson.M{"_id": i.Module})
	if err != nil {
		return module, err
	}

	if err = result.Decode(&module); err != nil {
		return module, err
	}

	err = module.AfterFind()
	if err != nil {
		return module, err
	}

	return module, nil
}

func (i *Instance) Create() (err error) {

	if err = i.BeforeCreate(); err != nil {
		return err
	}

	result, err := i.db().InsertOne(context.Background(), i)
	if err != nil {
		return err
	}

	id := result.InsertedID.(primitive.ObjectID)
	i.Id = id

	return nil
}

func (i *Instance) AfterFind() error {
	return nil
}

func (i *Instance) BeforeCreate() error {
	return nil
}

func RouteInstances(router chi.Router) {
	router.Post("/", createInstance)
	router.Get("/{id}", findInstance)
	router.Get("/{id}/func/{function}", runFunction)
	router.Get("/", findInstances)
}

func runFunction(w http.ResponseWriter, r *http.Request) {
	req, ctx, db := server.NewRequest(w, r, "instances")

	id, err := req.ParamObjectId("id")
	if err != nil {
		return
	}

	result := db.FindOne(ctx, bson.M{"_id": id})
	if err != nil {
		fmt.Println(err.Error())
	}

	var instance Instance
	err = result.Decode(&instance)
	if err != nil {
		return
	}

	result = server.From("modules").FindOne(ctx, bson.M{"_id": instance.Module})
	if err != nil {
		fmt.Println(err.Error())
	}

	var module Module
	err = result.Decode(&module)
	if err != nil {
		return
	}

	mod, err := module.Initialize(instance.Id.Hex())
	if err != nil {
		return
	}

	function := req.Param("function")
	if err != nil {

		return
	}

	run, err := mod.Run(function)
	if err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	req.ResolveRaw(run, http.StatusOK)
}

func findInstance(w http.ResponseWriter, r *http.Request) {
	req, ctx, db := server.NewRequest(w, r, "instances")

	id, err := req.ParamObjectId("id")
	if err != nil {
		return
	}

	result := db.FindOne(ctx, bson.M{"_id": id})
	if err != nil {
		fmt.Println(err.Error())
	}

	var model Instance
	err = result.Decode(&model)
	if err != nil {
		return
	}

	req.Resolve(model, http.StatusOK)
}

func createInstance(w http.ResponseWriter, r *http.Request) {
	req, _, _ := server.NewRequest(w, r, "instances")

	var model Instance
	req.DecodeModel(&model)

	err := model.Create()
	if err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	req.Resolve(model, http.StatusOK)
}

func findInstances(w http.ResponseWriter, r *http.Request) {
	req, ctx, db := server.NewRequest(w, r, "instances")

	var models []Instance
	cursor, err := db.Find(ctx, bson.M{})
	if err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	for cursor.Next(ctx) {
		var instance Instance
		if err = cursor.Decode(&instance); err != nil {
			return
		}

		if err = instance.AfterFind(); err != nil {
			return
		}

		models = append(models, instance)
	}

	req.Resolve(models, http.StatusOK)
}
