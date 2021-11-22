// Copyright (c) 2021 Braden Nicholson

package endpoint

import (
	"fmt"
	"udap/udap"
	"udap/udap/db"
)

type Grant struct {
	udap.Persistent
	EndpointId string `json:"endpointId"`
	InstanceId string `json:"instanceId"`
}

// Grant converts requestBodies into their associate pairs
func (e *Endpoint) Grant(body requestBody) (err error) {
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
	grant := Grant{
		EndpointId: e.Id,
		InstanceId: instanceId,
	}
	err = db.DB.Model(&Grant{}).Create(&grant).Error
	if err != nil {
		return err
	}
	udap.Log("Granted access from endpoint %s to instance %s", e.Id, instanceId)
	return err
}

// Revoke handles revoke requests
func (e *Endpoint) Revoke(body requestBody) (err error) {
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
	err = db.DB.Delete(&Grant{}, "endpoint_id = ? AND instance_id = ?", e.Id, instanceId).Error
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
	err := db.DB.Where(&Grant{}, "endpoint_id = ? AND instance_id = ?", e.Id,
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
	err = db.DB.Where(&Grant{}, "endpoint_id = ?", e.Id).First(&grants).Error
	if err != nil {
		return nil, err
	}
	for _, grant := range grants {
		instanceIds = append(instanceIds, grant.InstanceId)
	}
	return instanceIds, err
}

type Subscription struct {
	udap.Persistent
	EndpointId string `json:"endpointId"`
	InstanceId string `json:"instanceId"`
}

// subscribed determines if the endpoint is subscribed to an instance with the provided id
func (e *Endpoint) subscribed(instanceId string) (status bool) {
	// Initialize a cache variable to hold the database model
	var subscription Subscription
	// Attempt to find the first record
	err := db.DB.First(&subscription, "endpoint_id = ? AND instance_id = ?", e.Id, instanceId).Error
	if err != nil {
		// Return false if some error occurs
		return false
	}
	// Return true if the model is not empty
	return subscription != Subscription{}
}

// Subscribe is the high level WebSocket interface
func (e *Endpoint) Subscribe(body requestBody) (err error) {
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
	if err = db.DB.Create(&subscription).Error; err != nil {
		return err
	}
	// Return no errors
	return nil
}

// Unsubscribe is removes a subscription relation
func (e *Endpoint) Unsubscribe(body requestBody) (err error) {
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
	err = db.DB.Delete(&Subscription{}, "endpoint_id = ? AND instance_id = ?", e.Id, instanceId).Error
	if err != nil {
		return err
	}
	// Return no errors
	return nil
}
