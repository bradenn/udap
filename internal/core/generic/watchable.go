// Copyright (c) 2022 Braden Nicholson

package generic

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
)

type Identifiable interface {
	GetId() string
}

type Watchable[T Identifiable] struct {
	channel chan<- domain.Mutation
}

func (w *Watchable[T]) Emit(element T) error {
	if w.channel == nil {
		return fmt.Errorf("channel is null")
	}
	classification := strings.ToLower(reflect.TypeOf(element).Name())

	eId := element.GetId()

	payload := domain.Mutation{
		Status:    "update",
		Operation: classification,
		Body:      element,
		Id:        eId,
	}

	// Set a timer to cancel sending after 100 milliseconds
	timer := time.NewTimer(time.Millisecond * 500)
	select {
	// Attempt to push the payload to the channel
	case w.channel <- payload:
		// Cancel the timer if payload is sent
		timer.Stop()
		// Exit normally
		return nil
	case <-timer.C:
		log.Event("emit failed for '%s'", classification)
		// Exit quietly if the payload could not be sent
		return nil
	}

}

func (w *Watchable[T]) Watch(ref chan<- domain.Mutation) error {
	w.channel = ref
	return nil
}
