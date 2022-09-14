// Copyright (c) 2022 Braden Nicholson

package generic

import (
	"testing"
	"udap/internal/core/domain"
)

type Basic struct {
	Variable string
}

func (b Basic) GetId() string {
	return "thisIsAUniqueId"
}

type BasicService struct {
	Watchable[Basic]
}

func TestWatchable_Emit(t *testing.T) {
	bs := BasicService{}

	recv := make(chan domain.Mutation)

	err := bs.Watch(recv)
	if err != nil {
		t.Error(err)
	}

	toEmit := Basic{Variable: "testString"}

	go func() {
		err = bs.Emit(toEmit)
		if err != nil {
			t.Error(err)
		}
	}()

	r := <-recv
	if r.Id != toEmit.GetId() {
		t.Errorf("id does not match")
	}

	if r.Status != "update" {
		t.Errorf("unimplemented status indicator")
	}

	if r.Operation != "basic" {
		t.Errorf("class identification failed")
	}

}
