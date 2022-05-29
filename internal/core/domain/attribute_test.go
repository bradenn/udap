// Copyright (c) 2022 Braden Nicholson

package domain

import (
	"testing"
)

func TestAttribute_Sanity(t *testing.T) {
	a := Attribute{
		Value: "applesauce",
	}
	if a.Value != "applesauce" {
		t.Errorf("Value not set")
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
