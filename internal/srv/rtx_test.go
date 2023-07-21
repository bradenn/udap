// Copyright (c) 2023 Braden Nicholson

package srv

import (
	"testing"
	"udap/internal/controller"
)

func TestNewRtx(t *testing.T) {
	srv := &Server{}
	ctrl := &controller.Controller{}
	rtx := NewRtx(srv, ctrl, nil)
	rtx.Loaded()
}
