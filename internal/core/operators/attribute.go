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
	a.mutex.Lock()
	a.hooks[attribute.Id] = attribute.Channel
	a.mutex.Unlock()
	return nil
}

func (a *attributeOperator) Request(attribute *domain.Attribute, s string) error {
	attribute.Requested = time.Now()
	attribute.Request = s

	if attribute.Request != attribute.Value || time.Since(attribute.UpdatedAt) >= time.Millisecond*1000 {
		var channel chan domain.Attribute

		a.mutex.Lock()
		channel = a.hooks[attribute.Id]
		a.mutex.Unlock()

		if channel == nil {
			return fmt.Errorf("channel is not set")
		}

		channel <- *attribute
		attribute.Requested = time.Now()

		err := a.Set(attribute, s)
		if err != nil {
			return err
		}
	}
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
