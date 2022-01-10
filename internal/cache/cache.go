// Copyright (c) 2021 Braden Nicholson

package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
)

var Mem Cache
var memCtx context.Context

func init() {
	memCtx = context.Background()
}

func PutLn(value any, path ...string) error {
	err := Mem.Set(memCtx, strings.ToLower(strings.Join(path, ".")), value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetLn(path ...string) (interface{}, error) {
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
