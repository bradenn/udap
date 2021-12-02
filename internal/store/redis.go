// Copyright (c) 2021 Braden Nicholson

package store

import (
	"context"
	"fmt"
	"strings"
	"udap/internal/log"
)

var ctx context.Context

func init() {
	ctx = context.Background()

}

func WatchFn(key string, fn func(string) error) {
	ps := client.Subscribe(context.Background(), key)
	_, err := ps.Receive(context.Background())
	if err != nil {
		return
	}
	ch := ps.Channel()
	log.Log("SUBSCRIBED %s", key)
	go func() {
		for message := range ch {
			fmt.Println(message.Payload)
			err = fn(message.Payload)
			if err != nil {
				return
			}
		}
	}()

	if err != nil {
		return
	}
}

func Put(key string, value string) error {
	// Publish a message.
	err := client.Publish(ctx, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func AllOfType(key string) (string, error) {
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		return value, err
	}
	return value, nil
}

func PutLn(value string, path ...string) error {
	err := client.Set(ctx, strings.Join(path, "."), value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func Get(key string) (string, error) {
	value, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return value, err
	}
	return value, nil
}

func GetLn(path ...string) (string, error) {
	value, err := client.Get(ctx, strings.Join(path, ".")).Result()

	if err != nil {
		fmt.Println(err)
		return value, err
	}
	return value, nil
}

func Path(path ...string) (res []string, err error) {
	result, err := client.Keys(context.Background(), "modules.*").Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}
