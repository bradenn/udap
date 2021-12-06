// Copyright (c) 2021 Braden Nicholson

package models

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"math/rand"
	"time"
	"udap/internal/log"
	"udap/internal/store"
)

// TODO: Making endpoint mutable

type RequestBody struct {
	Name       string `json:"name"`
	TargetId   string `json:"targetId"`
	InstanceId string `json:"instanceId"`
}

// Endpoint represents a client device connected to the UDAP network
type Endpoint struct {
	store.Persistent
	// Name identifies the terminal
	Name string `json:"name" gorm:"unique"`
	// Type The specification of endpoint
	Type string `json:"type"`
	// Conn is the websocket connection with an endpoint
	Conn *websocket.Conn `gorm:"-" json:"-"`
	// key is used to identify new endpoints
	key string
	// registered defines the state of the endpoint
	registered bool
	// enrolledSince is the time the endpoint was enrolled
	enrolledSince time.Time `gorm:"-"`
}

// BeforeCreate is a hook function from gorm, called when an endpoint is inserted
func (e *Endpoint) BeforeCreate(_ *gorm.DB) error {
	e.key = randomSequence()
	return nil
}

// randomSequence generates a random id for use as a key
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

// Fetch populates the struct from the database
func (e *Endpoint) Fetch() error {
	err := store.DB.Where("id = ?", e.Id).First(e).Error
	if err != nil {
		return err
	}
	return nil
}

// AfterFind is a gorm hook, called when an endpoint is found
func (e *Endpoint) AfterFind(_ *gorm.DB) error {
	e.registered = false
	// endpoints[e.Id] = e
	return nil
}

// Enrollment

// Enroll adds an endpoint's websocket Conn to the active update queue
func (e *Endpoint) Enroll() (err error) {
	// If the endpoint is already enrolled, unenroll
	if e.enrolled() {
		err = e.Unenroll()
		if err != nil {
			return err
		}
	}
	// Run the socketCloseHandler when the endpoint is disconnected
	e.Conn.SetCloseHandler(e.socketCloseHandler())
	// Log the new state
	log.Sherlock("Endpoint '%s' connected.", e.Name)
	e.registered = true
	e.enrolledSince = time.Now()
	// Return no errors
	return nil
}

// Unenroll removes an endpoint from the active update queue
func (e *Endpoint) Unenroll() error {
	e.registered = false
	log.Sherlock("Endpoint '%s' disconnected after %s.", e.Name, time.Since(e.enrolledSince).Round(time.
		Millisecond*10))
	e.enrolledSince = time.Now()
	return nil
}

// enrolled determines whether the endpoint believes itself to be enrolled
func (e *Endpoint) enrolled() bool {
	return e.registered
}

// socketCloseHandler is called when a websocket Conn is terminated
func (e *Endpoint) socketCloseHandler() func(code int, text string) error {
	return func(code int, text string) error {
		err := e.Unenroll()
		if err != nil {
			return err
		}
		return nil
	}
}

// Grants

// Grant provides access of an endpoint to an instance
type Grant struct {
	store.Persistent
	EndpointId string `json:"endpointId"`
	InstanceId string `json:"instanceId"`
}

// Grant converts requestBodies into their associate pairs
func (e *Endpoint) Grant(body RequestBody) (err error) {
	// TODO: Authenticate User
	// Check if the endpoint has already been granted this instance
	if e.granted(body.InstanceId) {
		return fmt.Errorf("already granted")
	}
	// Attempt to grant the instance to the endpoint
	if err = e.grant(body.InstanceId); err != nil {
		return err
	}
	return nil
}

// grant creates grant records, permitting an endpoint to access an instance
func (e *Endpoint) grant(instanceId string) (err error) {
	// Allocate the grant structure
	grant := Grant{
		EndpointId: e.Id,
		InstanceId: instanceId,
	}
	// Attempt to insert the record
	if err = store.DB.Model(&Grant{}).Create(&grant).Error; err != nil {
		return err
	}
	// Return no errors
	return nil
}

// Revoke handles revoke requests
func (e *Endpoint) Revoke(body RequestBody) (err error) {
	// TODO: Authenticate User
	// Make sure the endpoint has the instance first
	if !e.granted(body.InstanceId) {
		return fmt.Errorf("endpoint has not been granted this instance")
	}
	// If the endpoint is concurrently subscribed to an instance, we need to unsubscribe first
	if e.subscribed(body.InstanceId) {
		// Attempt to unsubscribe from the instance
		err = e.unsubscribe(body.InstanceId)
		if err != nil {
			return err
		}
	}
	// Attempt to revoke the instance
	if err = e.revoke(body.InstanceId); err != nil {
		return err
	}
	// Return no errors
	return nil
}

// revoke removes an instance from the endpoints granted instances
func (e *Endpoint) revoke(instanceId string) (err error) {
	// Attempt to delete the database record for the grant
	err = store.DB.Delete(&Grant{}, "endpoint_id = ? AND instance_id = ?", e.Id, instanceId).Error
	if err != nil {
		return err
	}
	// Return no errors
	return nil
}

// granted determines where the instance has been granted the instance with the provided id
func (e *Endpoint) granted(instanceId string) (status bool) {
	// Initialize a cache variable to hold the database model
	var grant Grant
	// Attempt to find the first record
	err := store.DB.Where(&Grant{}, "endpoint_id = ? AND instance_id = ?", e.Id,
		instanceId).First(&grant).Error
	if err != nil {
		// Return false if some error occurs
		return false
	}
	// Return true if the model is not empty
	return grant != Grant{}
}

// grants returns an array of granted instance ids
func (e *Endpoint) grants() (instanceIds []string, err error) {
	var grants []Grant
	err = store.DB.Where(&Grant{}, "endpoint_id = ?", e.Id).First(&grants).Error
	if err != nil {
		return nil, err
	}
	for _, grant := range grants {
		instanceIds = append(instanceIds, grant.InstanceId)
	}
	return instanceIds, err
}

// Subscription

// Subscription represents a subscription of an endpoint to an instance
type Subscription struct {
	store.Persistent
	EndpointId string `json:"endpointId"`
	InstanceId string `json:"instanceId"`
}

func (s Subscription) TableName() string {
	return "subscriptions"
}

// Subscribe is the high level WebSocket interface
func (e *Endpoint) Subscribe(body RequestBody) (err error) {
	// Ensure the endpoint has been granted the instance
	if !e.granted(body.InstanceId) {
		return fmt.Errorf("instance not granted to endpoint")
	}
	// Make sure the endpoint is not already subscribed to the instance
	if e.subscribed(body.InstanceId) {
		return fmt.Errorf("already subscribed")
	}
	// Attempt to subscribe to the instance
	if err = e.subscribe(body.InstanceId); err != nil {
		return fmt.Errorf("internal error")
	}

	// store.WatchFn(fmt.Sprintf("instance.%s.buffer", body.InstanceId), e.onIns)
	// Return no error
	return nil
}

// subscribe is the low level module interface
func (e *Endpoint) subscribe(instanceId string) (err error) {
	// Initialize a subscription struct
	subscription := Subscription{
		EndpointId: e.Id,
		InstanceId: instanceId,
	}
	// Attempt to insert it into the subscriptions table
	if err = store.DB.Create(&subscription).Error; err != nil {
		return err
	}
	// Return no errors
	return nil
}

// Unsubscribe is removes a subscription relation
func (e *Endpoint) Unsubscribe(body RequestBody) (err error) {
	if !e.subscribed(body.InstanceId) {
		return fmt.Errorf("not subscribed")
	}
	if err = e.unsubscribe(body.InstanceId); err != nil {
		return fmt.Errorf("internal error")
	}
	return nil
}

// unsubscribe removes the database relation defining the subscription
func (e *Endpoint) unsubscribe(instanceId string) (err error) {
	// Attempt to delete the database record for the grant
	err = store.DB.Delete(&Subscription{}, "endpoint_id = ? AND instance_id = ?", e.Id, instanceId).Error
	if err != nil {
		return err
	}
	// Return no errors
	return nil
}

// subscribed determines if the endpoint is subscribed to an instance with the provided id
func (e *Endpoint) subscribed(instanceId string) (status bool) {
	// Initialize a cache variable to hold the database model
	var subscription Subscription
	// Attempt to find the first record
	err := store.DB.First(&subscription, "endpoint_id = ? AND instance_id = ?", e.Id, instanceId).Error
	if err != nil {
		// Return false if some error occurs
		return false
	}
	// Return true if the model is not empty
	return subscription != Subscription{}
}

// subscribed determines if the endpoint is subscribed to an instance with the provided id
func (e *Endpoint) subscriptions() (subscriptionIds []string, err error) {
	// Initialize a cache variable to hold the database model
	var subs []Subscription
	// Attempt to find the first record
	err = store.DB.Model(&Subscription{}).Find(&subs, "endpoint_id = ?", e.Id).Error
	if err != nil {
		// Return false if some error occurs
		return []string{}, nil
	}
	// Add each id to the ids list
	for _, sub := range subs {
		subscriptionIds = append(subscriptionIds, sub.InstanceId)
	}
	// Return true if the model is not empty
	return subscriptionIds, nil
}
