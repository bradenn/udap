// Copyright (c) 2021 Braden Nicholson

package cache

import (
	"strings"
	"udap/internal/log"
)

func WatchFn(key string, fn *func(string) error) {
	log.Sherlock("Watching: %s", strings.ToLower(key))
	ps := Mem.Subscribe(memCtx, strings.ToLower(key))
	_, err := ps.Receive(memCtx)
	if err != nil {
		log.Err(err)
	}
	ch := ps.Channel()
	go func(fn *func(string) error) {
		for message := range ch {
			err = (*fn)(message.Payload)
			if err != nil {
				log.Err(err)
			}
		}
	}(fn)
}

func Put(key string, value string) error {
	// Publish a message.
	err := Mem.Publish(memCtx, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func AllOfType(key string) (string, error) {
	value, err := Mem.Get(memCtx, key).Result()
	if err != nil {
		return value, err
	}
	return value, nil
}

func PutLn(value []byte, path ...string) error {
	err := Mem.Publish(memCtx, strings.ToLower(strings.Join(path, ".")), value).Err()
	if err != nil {
		log.Err(err)
	}
	return nil
}

func Get(key string) (string, error) {
	value, err := Mem.Get(memCtx, key).Result()
	if err != nil {
		return value, err
	}
	return value, nil
}

func GetLn(path ...string) (interface{}, error) {
	value, err := Mem.Get(memCtx, strings.ToLower(strings.Join(path, "."))).Result()
	if err != nil {
		return value, err
	}
	return value, nil
}

func Path(path ...string) (res []string, err error) {
	result, err := Mem.Keys(memCtx, "modules.*").Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
