package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"plugin"
	"udap/server"
	"udap/template"
)

type Module struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Environment string             `json:"environment"`
	Functions   []string           `json:"functions"`
	Path        string             `json:"path"`
	template.Metadata
}

// Module interface operations

func (m *Module) Initialize(instanceId string) (*template.Module, error) {
	module, err := m.load()
	if err != nil {
		return &template.Module{}, err
	}
	module.Configure([]byte(m.Environment), instanceId)
	return module, nil
}

func (m *Module) load() (*template.Module, error) {
	filename := fmt.Sprintf("./plugins/%s/%s.so", m.Path, m.Path)
	p, err := plugin.Open(filename)
	if err != nil {
		return &template.Module{}, err
	}
	lookup, err := p.Lookup("Export")
	if err != nil {
		return &template.Module{}, err
	}
	return lookup.(*template.Module), nil
}

func (m *Module) valid() bool {
	_, err := m.load()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// Hooks

func (m *Module) BeforeCreate() error {
	if !m.valid() {
		return fmt.Errorf("module does not exist")
	}
	module, err := m.load()
	if err != nil {
		return err
	}
	m.Metadata = module.Metadata()
	return nil
}

func (m *Module) AfterFind() error {
	load, err := m.load()
	if err != nil {
		return err
	}
	m.Functions = load.Functions()
	return nil
}

// API

func (m *Module) Route(router chi.Router) {
	router.Post("/", createModule)
	router.Get("/{id}", findModule)
	router.Get("/", findModules)
}

func findModule(w http.ResponseWriter, r *http.Request) {
	req, ctx, db := server.NewRequest(w, r, "modules")

	id, err := req.ParamObjectId("id")
	if err != nil {
		return
	}

	var model Module
	err = db.FindOne(ctx, bson.M{"_id": id}).Decode(&model)
	if err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	err = model.AfterFind()
	if err != nil {
		return
	}

	req.Resolve(model, http.StatusOK)
}

func createModule(w http.ResponseWriter, r *http.Request) {
	req, ctx, db := server.NewRequest(w, r, "modules")

	var model Module
	req.DecodeModel(&model)

	err := model.BeforeCreate()
	if err != nil {
		req.Reject(err.Error(), http.StatusConflict)
		return
	}

	result, err := db.InsertOne(ctx, &model)
	if err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	model.Id = result.InsertedID.(primitive.ObjectID)

	err = model.AfterFind()
	if err != nil {
		return
	}

	req.Resolve(model, http.StatusOK)
}

func findModules(w http.ResponseWriter, r *http.Request) {
	req, ctx, db := server.NewRequest(w, r, "modules")

	var model []Module

	results, err := db.Find(ctx, bson.D{})
	if err != nil {
		req.Reject(err.Error(), http.StatusConflict)
	}

	for results.Next(ctx) {
		mod := Module{}
		err = results.Decode(&mod)
		if err != nil {
			return
		}
		model = append(model, mod)
	}

	req.Resolve(model, http.StatusOK)
}
