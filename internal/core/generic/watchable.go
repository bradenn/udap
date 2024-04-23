// Copyright (c) 2022 Braden Nicholson

package generic

import (
	"reflect"
	"strings"
	"sync"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
)

type Identifiable interface {
	GetId() string
}

type Watchable[T Identifiable] struct {
	channels []chan<- domain.Mutation
	channel  chan<- domain.Mutation
	mutex    sync.RWMutex
}

func (w *Watchable[T]) Emit(element T) error {

	w.mutex.Lock()
	if w.channels == nil {
		w.channels = make([]chan<- domain.Mutation, 0)
	}
	w.mutex.Unlock()
	classification := strings.ToLower(reflect.TypeOf(element).Name())

	eId := element.GetId()

	payload := domain.Mutation{
		Status:    "update",
		Operation: classification,
		Body:      element,
		Id:        eId,
	}

	w.mutex.RLock()
	for _, channel := range w.channels {

		// Set a timer to cancel sending after 100 milliseconds
		timer := time.NewTimer(time.Millisecond * 500)
		select {
		// Attempt to push the payload to the channel
		case channel <- payload:
			// Cancel the timer if payload is sent
			timer.Stop()
			continue
		case <-timer.C:
			log.Event("emit failed for '%s'", payload)
			continue
		}
	}
	w.mutex.RUnlock()
	return nil
}

func (w *Watchable[T]) Watch(ref chan<- domain.Mutation) {

	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.channels == nil {
		w.channels = make([]chan<- domain.Mutation, 0)
	}

	w.channels = append(w.channels, ref)
}
