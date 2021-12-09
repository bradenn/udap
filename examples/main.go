// Copyright (c) 2021 Braden Nicholson

package examples

import "fmt"

type Entity interface {
	Name() string
	Module() string
	Path() string
	Style() string
	State() map[string]interface{}
}

func NewSwitch() Entity {
	return Default{
		module: "test",
		name:   "apple",
	}
}

type Default struct {
	module string
	name   string
	style  string
}

func (d Default) Style() string {
	return "default"
}

func (d Default) Name() string {
	return d.name
}

func (d Default) Module() string {
	return d.module
}

func (d Default) Path() string {
	return fmt.Sprintf("%s.%s", d.module, d.name)
}

func (d Default) State() map[string]interface{} {
	// TODO implement me
	panic("implement me")
}
