// Copyright (c) 2022 Braden Nicholson

package models

import (
	"testing"
	"time"
)

func TestAttribute_Sanity(t *testing.T) {
	a := Attribute{
		Value: "applesauce",
	}
	if a.Value != "applesauce" {
		t.Errorf("Value not set")
	}
}

func TestAttribute_Path(t *testing.T) {
	a := Attribute{
		Entity: "abc123",
		Key:    "cba321",
	}
	if a.Path() != "abc123.cba321" {
		t.Errorf("Attribute path is malformed")
	}
}

func TestAttribute_FnGet(t *testing.T) {
	a := Attribute{
		Entity: "abc123",
		Key:    "cba321",
	}

	a.FnGet(func() (string, error) {
		return "xyz123", nil
	})
	get, err := a.get()
	if err != nil {
		t.Error(err)
	}
	if get != "xyz123" {
		t.Error("function did not change value")
	}
}

func TestAttribute_FnPut(t *testing.T) {
	a := Attribute{
		Entity: "abc123",
		Key:    "cba321",
	}

	a.FnPut(func(value string) error {
		if value != "xyz123" {
			t.Error("function did not change value")
		}
		return nil
	})
	err := a.put("xyz123")
	if err != nil {
		t.Error(err)
	}

}

func TestAttribute_UpdateValue(t *testing.T) {
	a := Attribute{
		Request:   "applesauce",
		Requested: time.Now(),
	}
	err := a.UpdateValue("orangejuice", time.Now())
	if err == nil {
		t.Errorf("value should not be updated")
	}
	if a.Request != "applesauce" {
		t.Errorf("value should not be overwritten")
	}
}

func TestAttribute_SendRequest(t *testing.T) {
	a := Attribute{
		Request:   "applesauce",
		Requested: time.Now(),
	}
	err := a.SendRequest("orangejuice")
	if err == nil {
		t.Errorf("functions not set, should return error")
	}
}

func TestAttribute_AsInt(t *testing.T) {
	a := Attribute{
		Value: "100",
	}
	if a.AsInt() != 100 {
		t.Errorf("failed to convert string to int")
	}
}

func TestAttribute_AsFloat(t *testing.T) {
	a := Attribute{
		Value: "123.456",
	}
	if a.AsFloat() != 123.456 {
		t.Errorf("failed to convert string to float")
	}
}

func TestAttribute_AsBool(t *testing.T) {
	a := Attribute{
		Value: "false",
	}
	if a.AsBool() != false {
		t.Errorf("failed to convert string to bool")
	}
}
