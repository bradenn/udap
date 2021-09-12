package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"net/http"
	"plugin"
	"udap/server"
	"udap/template"
)

type Module struct {
	Persistent
	template.Metadata
	Environment string   `json:"environment"`
	Functions   []string `json:"functions" gorm:"-"`
	Path        string   `json:"path"`
}

func (m *Module) valid() bool {
	_, err := m.load()
	if err != nil {
		return false
	}
	return true
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

func (m *Module) Initialize(instanceId uuid.UUID) (*template.Module, error) {
	module, err := m.load()
	if err != nil {
		return &template.Module{}, err
	}
	module.Configure([]byte(m.Environment), instanceId)
	return module, nil
}

func (m *Module) AfterFind(_ *gorm.DB) error {
	load, err := m.load()
	if err != nil {
		return err
	}
	m.Functions = load.Functions()
	return nil
}
func (m *Module) BeforeCreate(_ *gorm.DB) error {
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

func (m *Module) Route(router chi.Router) {
	router.Post("/", createModule)
	router.Get("/{id}", findModule)
	router.Get("/", findModules)
}

func findModule(w http.ResponseWriter, r *http.Request) {
	req, db := server.NewRequest(w, r)

	var model Module

	id := req.Param("id")

	err := db.Where("id = ?", id).First(&model).Error
	if err != nil {
		req.Reject(err.Error(), http.StatusNotFound)
		return
	}

	req.Resolve(model, http.StatusOK)
}

func createModule(w http.ResponseWriter, r *http.Request) {
	req, db := server.NewRequest(w, r)

	var model Module

	req.DecodeModel(&model)

	err := db.Create(&model).Error
	if err != nil {
		req.Reject(err.Error(), http.StatusConflict)
		return
	}

	req.Resolve(model, http.StatusOK)
}

func findModules(w http.ResponseWriter, r *http.Request) {
	req, db := server.NewRequest(w, r)

	var model []Module

	err := db.Model(&model).Find(&model).Error
	if err != nil {
		req.Reject(err.Error(), http.StatusConflict)
		return
	}

	req.Resolve(model, http.StatusOK)
}
