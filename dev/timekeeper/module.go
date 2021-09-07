package main

import (
	"time"
)

type Function func(...interface{}) interface{}

type module struct {
	functions map[string]Function
	name      string
	desc      string
}

func (m *module) Functions() (functions []string) {
	for s := range m.functions {
		functions = append(functions, s)
	}
	return functions
}

func (m *module) emplace(name string, function Function) {
	m.functions[name] = function
}

func (m module) Run(name string, payload ...interface{}) interface{} {
	return m.functions[name](payload)
}

var Module module

func init() {

	Module = module{
		functions: map[string]Function{},
	}

	Module.emplace("timekeeper.local", GetLocalTime)
}

func GetLocalTime(b ...interface{}) interface{} {
	return time.Now().String()
}
