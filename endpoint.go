package main

import (
	"github.com/go-chi/chi"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"time"
	"udap/server"
)

// Endpoint represents a client device connected to the UDAP network
type Endpoint struct {
	Persistent
	// Name identifies the terminal
	Name string `json:"name" gorm:"unique"`
	// Instances define what modules the endpoint can interact with
	Instances []Instance `json:"instances" gorm:"many2many:endpointInstances;"`
	// Enabled must be true for the endpoint to be functional
	Enabled bool `json:"enabled"`
	// Key is used to identify new endpoints
	Key string `json:"key"`
}

func randomSequence() string {
	template := "abcdefghijklmnopqrstuvwxyz"
	var out string
	rand.Seed(time.Now().Unix())
	for i := 0; i < 8; i++ {
		r := rand.Intn(26)
		u := template[r]
		out += string(u)
	}
	return out
}

func (e *Endpoint) BeforeCreate(_ *gorm.DB) error {
	e.Enabled = false
	e.Key = randomSequence()
	return nil
}

func (e *Endpoint) AfterFind(_ *gorm.DB) error {
	return nil
}

func (e *Endpoint) GetInstance(instanceId string, db *gorm.DB) (Instance, error) {
	instance := Instance{}
	err := db.Model(&Instance{}).Preload("Module").Where("id = ?", instanceId).First(&instance).Error
	if err != nil {
		return instance, err
	}
	return instance, nil
}

func (e *Endpoint) Find(db *gorm.DB) error {
	err := db.Model(e).Where("id = ?", e.Id.String()).First(e).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) AddInstance(instance Instance) {
	e.Instances = append(e.Instances, instance)
}

func (e *Endpoint) Save(db *gorm.DB) error {
	err := db.Model(e).Where("id = ?", e.Id.String()).Save(e).Error
	if err != nil {
		return err
	}
	return nil
}

func RouteEndpoint(router chi.Router) {

	router.Get("/register/{accessKey}", registerEndpoint)

}

func registerEndpoint(w http.ResponseWriter, r *http.Request) {
	req, db := server.NewRequest(w, r)
	key := req.Param("accessKey")

	endpoint := Endpoint{}

	err := db.Model(&Endpoint{}).Where("key = ?", key).First(&endpoint).Error
	if err != nil {
		req.Reject(err, http.StatusBadRequest)
	}

	jwt, err := server.SignUUID(endpoint.Id.String())
	if err != nil {
		req.Reject("Internal Error", http.StatusInternalServerError)
		return
	}

	resolve := map[string]interface{}{"token": jwt}

	req.Resolve(resolve, http.StatusOK)
}

// func createEndpoint(writer http.ResponseWriter, request *http.Request) {
// 	req, db := server.NewRequest(writer, request)
//
// 	var err error
// 	var model Endpoint
//
// 	req.DecodeModel(&model)
// 	db.Create(&model)
//
// 	if err := db.Error; err != nil {
// 		req.Reject(err.Error(), http.StatusNotFound)
// 		return
// 	}
//
// 	jwt, err := server.SignUUID(model.Id.String())
// 	if err != nil {
// 		req.Reject("Internal Error", http.StatusInternalServerError)
// 		return
// 	}
//
// 	resolve := map[string]interface{}{"token": jwt}
//
// 	req.Resolve(resolve, http.StatusOK)
// }
//
// func findEndpoints(writer http.ResponseWriter, request *http.Request) {
// 	req, db := server.NewRequest(writer, request)
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
//
// func findEndpoint(writer http.ResponseWriter, request *http.Request) {
// 	req, db := server.NewRequest(writer, request)
//
// 	id := req.Param("id")
//
// 	var model Endpoint
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
