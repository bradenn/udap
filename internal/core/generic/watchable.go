// Copyright (c) 2022 Braden Nicholson

package generic

import (
	"fmt"
	"udap/internal/core/domain"
	"udap/internal/core/domain/common"
)

type Watchable struct {
	channel chan<- domain.Mutation
}

func (w *Watchable) emit(element *common.Persistent) error {
	if w.channel == nil {
		return fmt.Errorf("channel is null")
	}

	w.channel <- domain.Mutation{
		Status:    "update",
		Operation: "attribute",
		Body:      *element,
		Id:        element.Id,
	}
	return nil
}

func (w *Watchable) Watch(ref chan<- domain.Mutation) error {
	if w.channel != nil {
		return fmt.Errorf("channel in use")
	}
	w.channel = ref
	return nil
}
