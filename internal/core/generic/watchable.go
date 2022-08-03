// Copyright (c) 2022 Braden Nicholson

package generic

import (
	"fmt"
	"time"
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
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

func (w *Watchable) emit(element *common.Persistent) error {
	if w.channel == nil {
		return fmt.Errorf("channel is null")
	}

	payload := domain.Mutation{
		Status:    "update",
		Operation: w.classification,
		Body:      *element,
		Id:        element.Id,
	}
	// Set a timer to cancel sending after 50 milliseconds
	timer := time.NewTimer(time.Millisecond * 50)
	select {
	// Attempt to push the payload to the channel
	case w.channel <- payload:
		// Cancel the timer if payload is sent
		timer.Stop()
		// Exit normally
		return nil
	case <-timer.C:
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
