// Copyright (c) 2022 Braden Nicholson

package operators

import (
	"fmt"
	"sync"
	"time"
	"udap/internal/core/domain"
	"udap/internal/core/ports"
)

type attributeOperator struct {
	hooks   map[string]chan domain.Attribute
	request chan domain.Attribute
	mutex   sync.RWMutex
}

func NewAttributeOperator() ports.AttributeOperator {
	return &attributeOperator{
		hooks:   map[string]chan domain.Attribute{},
		mutex:   sync.RWMutex{},
		request: make(chan domain.Attribute, 8),
	}
}

func (a *attributeOperator) Register(attribute *domain.Attribute) error {
	if attribute.Id == "" {
		return fmt.Errorf("invalid attribute id")
	}

	a.hooks[attribute.Id] = attribute.Channel
	return nil
}

func (a *attributeOperator) Request(attribute *domain.Attribute, s string) error {
	attribute.Requested = time.Now()
	attribute.Request = s

	// Find the correct channel
	channel, ok := a.hooks[attribute.Id]
	if !ok {
		return fmt.Errorf("channel is not set")
	}

	// Send the request to the module for processing
	channel <- *attribute

	return nil
}

func (a *attributeOperator) Set(attribute *domain.Attribute, s string) error {
	// If the attribute handler is not set, return an error

	attribute.Request = s
	attribute.Value = s

	return nil
}

func (a *attributeOperator) Update(attribute *domain.Attribute, val string, stamp time.Time) error {
	//// If a request has been made in the last five seconds, and has been unresolved, ignore this update
	//if attribute.Requested.Before(stamp) && attribute.Request != val && time.Since(attribute.Requested) < 300*time.Millisecond {
	//	return nil
	//}
	// Set the value

	attribute.Updated = stamp
	err := a.Set(attribute, val)
	if err != nil {
		return err
	}
	// Return no errors
	return nil
}
