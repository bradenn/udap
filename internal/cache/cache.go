// Copyright (c) 2021 Braden Nicholson

package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
)

var Mem Cache
var memCtx context.Context

type Cache struct {
	*redis.Client
}

func (c *Cache) Name() (name string) {
	return "cache"
}

func (c *Cache) Dependency() (level int) {
	return 1
}

// Load will begin the main-sequence activities of the parent struct
func (c *Cache) Load() (err error) {
	c.Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	memCtx = context.Background()
	Mem = *c
	return nil
}

func (c *Cache) WatchFn(key string, fn func(string) error) {
	ps := c.Subscribe(context.Background(), strings.ToLower(key))
	_, err := ps.Receive(context.Background())
	if err != nil {
		return
	}
	ch := ps.Channel()
	go func() {
		for message := range ch {
			err = fn(message.String())
			if err != nil {
				return
			}
		}
	}()
}

func (c *Cache) Run(ctx context.Context) (err error) {

	return nil
}

func (c *Cache) Cleanup() (err error) {
	return nil
}
