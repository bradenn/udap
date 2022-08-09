// Copyright (c) 2022 Braden Nicholson

package generic

import (
	"fmt"
	"time"
	"udap/internal/core/domain"
	"udap/internal/log"
)

type Watchable struct {
	channel        chan<- domain.Mutation
	classification string
}

func NewWatchable(classification string) Watchable {
	return Watchable{
		channel:        nil,
		classification: classification,
	}
}

func (w *Watchable) Emit(element any, id string) error {
	if w.channel == nil {
		return fmt.Errorf("channel is null")
	}

	payload := domain.Mutation{
		Status:    "update",
		Operation: w.classification,
		Body:      element,
		Id:        id,
	}
	// Set a timer to cancel sending after 100 milliseconds
	timer := time.NewTimer(time.Millisecond * 100)
	select {
	// Attempt to push the payload to the channel
	case w.channel <- payload:
		// Cancel the timer if payload is sent
		timer.Stop()
		// Exit normally
		return nil
	case <-timer.C:
		log.Event("emit failed for '%s'", w.classification)
		// Exit quietly if the payload could not be sent
		return nil
	}

}

func (w *Watchable) Watch(ref chan<- domain.Mutation) error {
	if w.channel != nil {
		return fmt.Errorf("channel in use")
	}
	w.channel = ref
	return nil
}
