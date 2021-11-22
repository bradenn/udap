// Copyright (c) 2021 Braden Nicholson

package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

type Named interface {
	Namespace() string
	UUID() string
}

func Save(named Named) error {
	marshal, err := json.Marshal(named)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("%s.%s", named.Namespace(), named.UUID()))
	err = client.Set(context.Background(), fmt.Sprintf("%s.%s", named.Namespace(), named.UUID()), marshal, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func Put(key string, value string) error {
	err := client.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func PutLn(value string, path ...string) error {
	err := client.Set(context.Background(), strings.Join(path, "."), value, 0).Err()
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
	value, err := client.Get(context.Background(), strings.Join(path, ".")).Result()
	if err != nil {
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

func red() {

}
