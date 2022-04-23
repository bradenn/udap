// Copyright (c) 2022 Braden Nicholson

package models

import (
	"testing"
)

func TestPath(t *testing.T) {
	a := Attribute{
		Entity: "abc123",
		Key:    "cba321",
	}
	if a.Path() != "abc123.cba321" {
		t.Errorf("Attribute path is malformed")
	}
}
