// Copyright (c) 2022 Braden Nicholson

package attribute

import (
	"fmt"
	"time"
	"udap/internal/core/domain"
)

type attributeOperator struct {
	hooks map[string]chan domain.Attribute
}

func NewOperator() domain.AttributeOperator {
	return &attributeOperator{
		hooks: map[string]chan domain.Attribute{},
	}
}

func (a attributeOperator) Register(attribute *domain.Attribute) error {
	a.hooks[attribute.Id] = attribute.Channel
	return nil
}

func (a attributeOperator) Request(attribute *domain.Attribute, s string) error {
	err := a.Set(attribute, s)
	if err != nil {
		return err
	}
	return nil
}

func (a attributeOperator) Set(attribute *domain.Attribute, s string) error {
	// If the attribute handler is not set, return an error
	channel := a.hooks[attribute.Id]

	attribute.Request = s

	attribute.Value = s

	attribute.Requested = time.Now()

	if channel == nil {
		return fmt.Errorf("channel is not open")
	}

	channel <- *attribute

	return nil
}

func (a attributeOperator) Update(attribute *domain.Attribute, val string, stamp time.Time) error {
	// If a request has been made in the last five seconds, and has been unresolved, ignore this update
	if attribute.Requested.Before(stamp) && attribute.Request != val && time.Since(attribute.Requested) < 5*time.Second {
		return fmt.Errorf("OVERWRITES REQUEST")
	}
	// Update the request value (since the request can be external)
	attribute.Request = val
	// Set the value
	err := a.Set(attribute, val)
	if err != nil {
		return err
	}
	// Return no errors
	return nil
}
