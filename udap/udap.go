// Copyright (c) 2021 Braden Nicholson

package udap

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type StoreBuffer struct {
	data string
}

type Pair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Udap struct {
	data          chan Pair
	storageBuffer map[string]string
	redis         redis.Client
}

func Print(path string) {
	fmt.Println("#", udap.storageBuffer[path])
}

func (u *Udap) dataChannel() {
	for buffer := range u.data {
		fmt.Println(buffer)
	}
}

var udap *Udap
var ctx = context.Background()

func New() *Udap {
	if udap != nil {
		return udap
	}
	dataChannel := make(chan Pair)
	udap = &Udap{
		data:          dataChannel,
		storageBuffer: map[string]string{},
	}

	go func() {
		for s := range dataChannel {
			udap.storageBuffer[s.Key] = s.Value
		}
	}()
	return udap
}

func startup() {

}
