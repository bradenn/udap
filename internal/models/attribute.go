// Copyright (c) 2022 Braden Nicholson

package models

import (
	"fmt"
	"strconv"
	"time"
)

type FuncPut func(value string) error
type FuncGet func() (string, error)

type Attribute struct {
	Id      string    `json:"id"`
	Value   string    `json:"value"`
	Updated time.Time `json:"updated"`

	Request   string    `json:"request"`
	Requested time.Time `json:"requested"`
	Entity    string    `json:"entity"`
	Key       string    `json:"key"`
	Type      string    `json:"type"`
	Order     int       `json:"order"`
	put       FuncPut
	get       FuncGet
}

// Path Returns a unique identifier bound to an entity
func (a *Attribute) Path() string {
	return fmt.Sprintf("%s.%s", a.Entity, a.Key)
}

// SetValue overrides any existing value
func (a *Attribute) SetValue(val string) {
	a.Request = val
	// Overwrite the value
	a.Value = val
	// Update the timestamp for the current values time
	a.Updated = time.Now()
}

// UpdateValue attempts to write an update to the attribute
func (a *Attribute) UpdateValue(val string, stamp time.Time) error {
	// If a request has been made in the last five seconds, and has been unresolved, ignore this update
	if a.Requested.Before(stamp) && a.Request != val && time.Since(a.Requested) < 5*time.Second {
		return fmt.Errorf("OVERWRITES REQUEST")
	}
	// Update the request value (since the request can be external)
	a.Request = val
	// Set the value
	a.SetValue(val)
	// Return no errors
	return nil
}

// SendRequest attempts to send a change to the attribute handler
func (a *Attribute) SendRequest(val string) error {
	// If the attribute handler is not set, return an error
	if a.put == nil {
		return fmt.Errorf("attribute put function not connected")
	}
	// Register the request
	a.Request = val
	// Mark the request's time
	a.Requested = time.Now()
	// Attempt to send the value
	err := a.put(val)
	if err != nil {
		return err
	}
	// Set the value
	a.SetValue(val)
	// Return no errors
	return nil
}

// FnPut registers the attributes set function
func (a *Attribute) FnPut(put FuncPut) {
	a.put = put
}

// FnGet registers the attributes get function
func (a *Attribute) FnGet(get FuncGet) {
	a.get = get
}

func (a *Attribute) AsInt() int {
	parsed, err := strconv.ParseInt(a.Value, 10, 64)
	if err != nil {
		return 0
	}
	return int(parsed)
}

func (a *Attribute) AsFloat() float64 {
	parsed, err := strconv.ParseFloat(a.Value, 64)
	if err != nil {
		return 0.0
	}
	return parsed
}

func (a *Attribute) AsBool() bool {
	parsed, err := strconv.ParseBool(a.Value)
	if err != nil {
		return false
	}
	return parsed
}
