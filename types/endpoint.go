package types

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
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
	//
	Connection    *websocket.Conn `gorm:"-"`
	Subscriptions []string        `gorm:"-"`
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

func (e *Endpoint) GetInstance(instanceId string) (Instance, error) {
	instance := Instance{}
	err := db.Model(&Instance{}).Preload("Module").Where("id = ?", instanceId).First(&instance).Error
	if err != nil {
		return instance, err
	}
	return instance, nil
}

func (e *Endpoint) GrantInstance(instance Instance) error {
	for _, candidate := range e.Instances {
		if candidate.Id == instance.Id {
			return fmt.Errorf("instance already granted")
		}
	}
	e.Instances = append(e.Instances, instance)

	err := db.Model(e).Association("Instances").Append(&instance)
	if err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) RevokeInstance(instanceId string) (err error) {
	for i, instance := range e.Instances {
		if instance.Id.String() == instanceId {
			e.Instances = append(e.Instances[:i], e.Instances[i+1:]...)
			break
		}
	}
	for _, instance := range e.Instances {
		fmt.Println(instance.Name, ":", instance.Id.String(), instance.Module.Name, instance.Module.Id)

	}
	err = db.Model(&e).Where("id = ?", e.Id.String()).Updates(&e).Error
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
	err := db.Model(e).Where("id = ?", e.Id.String()).First(e).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) AddInstance(instance Instance) {
	e.Instances = append(e.Instances, instance)
}

func (e *Endpoint) Save() error {
	err := db.Model(e).Where("id = ?", e.Id.String()).Save(e).Error
	if err != nil {
		return err
	}
	return nil
}

func RouteEndpoint(router chi.Router) {
	router.Get("/register/{accessKey}", registerEndpoint)
	// router.Get("/socket/{jwt}", handleWebsockets)
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
