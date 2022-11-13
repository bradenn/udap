// Copyright (c) 2022 Braden Nicholson

package domain

type Mutation struct {
	Status    string `json:"status"`
	Operation string `json:"operation"`
	Body      any    `json:"body"`
	Id        string `json:"id"`
}

type Observer struct {
}

func (o *Observer) emit() {

}

type Observable interface {
	Watch(chan<- Mutation)
	EmitAll() error
}
