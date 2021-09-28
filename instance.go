package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Instance is a subclass of Module. Instance carries instance-related environment information.
type Instance struct {
	Persistent
	// Name refers to the name of the instance, not the name of the module.
	Name string `json:"name" gorm:"unique"`
	// Description briefly describe the nature of the instance, not the module.
	Description string `json:"description"`
	// Permission required to access the instance
	Permission string `json:"permission" bson:"permission"`
	// Module is the actual reference to the plugin module.
	Module   Module    `json:"module" gorm:"foreignKey:ModuleId"`
	ModuleId uuid.UUID `json:"moduleId"`
	// Config holds instance related environment information in JSON format.
	Config string `json:"config"`
	// Log refers to the STD output of the module
	Log string `json:"log"`
}

func (i *Instance) AfterFind(_ *gorm.DB) error {
	return nil
}

func (i *Instance) BeforeCreate(_ *gorm.DB) error {

	mod, err := i.Module.Initialize()
	if err != nil {
		return err
	}

	instance, err := mod.InitInstance()
	if err != nil {
		return err
	}

	i.Config = instance

	return nil
}

func (i *Instance) Poll() (string, error) {
	mod, err := i.Module.Initialize()
	if err != nil {
		return "", err
	}
	poll, err := mod.Poll(i.Config)
	if err != nil {
		return "", err
	}
	return poll, nil
}

func (i *Instance) Run(function string) (string, error) {
	mod, err := i.Module.Initialize()
	if err != nil {
		return "", err
	}

	result, err := mod.Run(i.Config, function)
	if err != nil {
		return "", err
	}

	return result, nil
}

//
//
// func runFunction(w http.ResponseWriter, r *http.Request) {
// 	req, db := server.NewRequest(w, r)
//
// 	id := req.Param("id")
//
// 	db.Where("id = ?", id)
// 	if err := db.Error; err != nil {
// 		req.Reject(err.Error(), http.StatusNotFound)
// 		return
// 	}
//
// 	req.ResolveRaw("", http.StatusOK)
// }
//
// func RouteInstances(router chi.Router) {
// 	router.Post("/", createInstance)
// 	router.Get("/{id}", findInstance)
// 	router.Get("/{id}/func/{function}", runFunction)
// 	router.Get("/", findInstances)
// }

// func findInstance(w http.ResponseWriter, r *http.Request) {
// 	req, db := server.NewRequest(w, r)
//
// 	id := req.Param("id")
//
// 	var model Instance
//
// 	db.Where("id = ?", id).First(&model)
//
// 	if err := db.Error; err != nil {
// 		req.Reject(err.Error(), http.StatusNotFound)
// 		return
// 	}
//
// 	req.Resolve(model, http.StatusOK)
// }
//
// func createInstance(w http.ResponseWriter, r *http.Request) {
// 	req, db := server.NewRequest(w, r)
//
// 	var model Instance
//
// 	req.DecodeModel(&model)
// 	db.Create(&model)
//
// 	if err := db.Error; err != nil {
// 		req.Reject(err.Error(), http.StatusNotFound)
// 		return
// 	}
//
// 	req.Resolve(model, http.StatusOK)
// }
//
// func findInstances(w http.ResponseWriter, r *http.Request) {
// 	req, db := server.NewRequest(w, r)
// 	var models []Endpoint
//
// 	db.Find(&models)
//
// 	if err := db.Error; err != nil {
// 		req.Reject(err.Error(), http.StatusNotFound)
// 		return
// 	}
//
// 	req.Resolve(models, http.StatusOK)
// }
