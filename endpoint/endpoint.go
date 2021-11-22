// Copyright (c) 2021 Braden Nicholson

package endpoint

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"udap/module"
	"udap/server"
	"udap/udap"
	"udap/udap/db"
	"udap/udap/store"
)

// TODO: Making endpoint mutable

type requestBody struct {
	Name       string `json:"name"`
	TargetId   string `json:"targetId"`
	InstanceId string `json:"instanceId"`
}

// Endpoint represents a client device connected to the UDAP network
type Endpoint struct {
	udap.Persistent
	// Name identifies the terminal
	Name string `json:"name" gorm:"unique"`
	// key is used to identify new endpoints
	key string

	// Instances define what modules the endpoint can interact with
	enrolled      bool               `gorm:"-"`
	connection    *websocket.Conn    `gorm:"-"`
	Instances     []*module.Instance `json:"instances" gorm:"many2many:endpointInstances;"`
	Subscriptions []*module.Instance `json:"subscriptions" gorm:"many2many:endpointSubscriptions;"`
}

func New(name string) (*Endpoint, error) {
	endpoint := &Endpoint{
		Name:          name,
		Instances:     nil,
		key:           "",
		enrolled:      false,
		connection:    nil,
		Subscriptions: nil,
	}
	err := db.DB.Model(&Endpoint{}).Create(&endpoint).Error
	if err != nil {
		return endpoint, err
	}
	return endpoint, nil
}

func Subscribe(instanceId string) {

}

func Unsubscribe(instanceId string) {

}

func (e *Endpoint) Namespace() string {
	return "endpoint"
}

func (e *Endpoint) Store() error {

	err := store.PutLn("endpoint.%s.name", strings.ToLower(e.Name))
	if err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) BeforeCreate(_ *gorm.DB) error {
	e.key = randomSequence()
	return nil
}

func (e *Endpoint) AfterFind(_ *gorm.DB) error {
	return nil
}

func (e *Endpoint) GetInstance(instanceId string) (module.Instance, error) {
	instance := module.Instance{}
	err := db.DB.Model(&module.Instance{}).Preload("Module").Where("id = ?", instanceId).First(&instance).Error
	if err != nil {
		return instance, err
	}
	return instance, nil
}

func (e *Endpoint) GrantInstance(instance module.Instance) error {
	for _, candidate := range e.Instances {
		if candidate.Id == instance.Id {
			return fmt.Errorf("instance already granted")
		}
	}
	e.Instances = append(e.Instances, &instance)

	err := db.DB.Model(e).Association("Instances").Append(&instance)
	if err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) RevokeInstance(instanceId string) (err error) {
	for i, instance := range e.Instances {
		if instance.Id == instanceId {
			e.Instances = append(e.Instances[:i], e.Instances[i+1:]...)
			break
		}
	}
	for _, instance := range e.Instances {
		fmt.Println(instance.Name, ":", instance.Id, instance.Module.Name, instance.Module.Id)

	}
	err = db.DB.Model(&e).Where("id = ?", e.Id).Updates(&e).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) DeleteInstance(instanceId string) error {
	err := e.RevokeInstance(instanceId)
	if err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) Find() error {
	err := db.DB.Model(e).Where("id = ?", e.Id).First(e).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) AddInstance(instance module.Instance) {
	e.Instances = append(e.Instances, &instance)
}

func (e *Endpoint) Save() error {
	err := db.DB.Model(e).Where("id = ?", e.Id).Save(e).Error
	if err != nil {
		return err
	}
	return nil
}

func RouteEndpoint(router chi.Router) {
	router.Get("/register/{accessKey}", registerEndpoint)
}

func registerEndpoint(w http.ResponseWriter, r *http.Request) {
	req, _ := server.NewRequest(w, r)
	key := chi.URLParam(r, "accessKey")
	endpoint := Endpoint{}

	err := db.DB.Model(&Endpoint{}).Where("key = ?", key).First(&endpoint).Error
	if err != nil {
		req.Reject(err, http.StatusBadRequest)
	}

	jwt, err := server.SignUUID(endpoint.Id)
	if err != nil {
		req.Reject("Internal Error", http.StatusInternalServerError)
		return
	}

	resolve := map[string]interface{}{"token": jwt}

	req.Resolve(resolve, http.StatusOK)
}
