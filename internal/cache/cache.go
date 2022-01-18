// Copyright (c) 2021 Braden Nicholson

package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
	"udap/internal/log"
)

var Mem Cache
var memCtx context.Context

func init() {
	memCtx = context.Background()
}

func WatchFn(key string, fn func(string) error) {
	log.Event("Watching: %s", strings.ToLower(key))
	ps := Mem.Subscribe(memCtx, strings.ToLower(key))
	ch := ps.Channel()
	go func(fn func(string) error) {
		for message := range ch {
			err := (fn)(message.Payload)
			if err != nil {
				log.Err(err)
			}
		}
	}(fn)
}

func Watch(key string, fn func(string) error) {
	log.Event("Watching: %s", strings.ToLower(key))
	ps := Mem.Subscribe(memCtx, strings.ToLower(key))
	ch := ps.Channel()
	go func(fn func(string) error) {
		for message := range ch {
			err := (fn)(message.Payload)
			if err != nil {
				log.Err(err)
			}
		}
	}(fn)
}

func PutLn(value any, path ...string) error {
	err := Mem.Publish(memCtx, strings.ToLower(strings.Join(path, ".")), value).Err()
	if err != nil {
		return err
	}
	err = Mem.Set(memCtx, strings.ToLower(strings.Join(path, ".")), value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func PutOne(value any, path ...string) error {
	err := Mem.Publish(memCtx, strings.ToLower(strings.Join(path, ".")), value).Err()
	if err != nil {
		return err
	}
	err = Mem.Set(memCtx, strings.ToLower(strings.Join(path, ".")), value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetLn(path ...string) (any, error) {
	value, err := Mem.Get(memCtx, strings.ToLower(strings.Join(path, "."))).Result()
	if err != nil {
		return value, err
	}
	return value, nil
}

type Cache struct {
	*redis.Client
}

func NewCache() (Cache, error) {

	Mem.Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return Mem, nil
}
