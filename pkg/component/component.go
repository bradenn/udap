// Copyright (c) 2021 Braden Nicholson

package component

import "sync"

type Runner struct {
	components sync.Map
}

func (r *Runner) Load(c ...Component) {
	for _, component := range c {
		err := component.Load()
		if err == nil {
			return
		}
		r.components.LoadOrStore(component.Id(), component)
	}

}

func (r *Runner) Run() {

}

type Component interface {
	Id() string
	Name() string
	Load() error
	Run() error
	Exit() error
}

type Default struct {
}

func (d Default) Id() string {
	return "nil"
}

func (d Default) Name() string {
	return "unset"
}

func (d Default) Load() error {
	return nil
}

func (d Default) Run() error {
	return nil
}

func (d Default) Exit() error {
	return nil
}

type Message struct {
	Id   string `json:"id"`   // module, daemon, agent, etc
	Type string `json:"type"` // module, daemon, agent, etc
	Body string `json:"body"`
}
