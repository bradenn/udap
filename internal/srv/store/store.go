// Copyright (c) 2024 Braden Nicholson

package store

import (
	"fmt"
	redistimeseries "github.com/RedisTimeSeries/redistimeseries-go"
	"os"
)

type Store struct {
	client *redistimeseries.Client
}

func NewStore() *Store {

	redisHost := os.Getenv("redisHost")
	redisPort := os.Getenv("redisPort")
	//redisPass := os.Getenv("redisPass")
	//
	//client := redis.NewClient(&redis.Options{
	//	Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
	//	Password: redisPass, // no password set
	//	DB:       0,         // use default DB
	//})

	var client = redistimeseries.NewClient(fmt.Sprintf("%s:%s", redisHost, redisPort), "", nil)

	return &Store{
		client: client,
	}
}

func (s *Store) Traces(from int64, to int64, window int, mode string, filter []string) (map[string]map[int64]float64, map[string]map[string]string, error) {

	options := redistimeseries.DefaultMultiRangeOptions

	options.SetAggregation(redistimeseries.AggregationType(mode), window).SetWithLabels(true).SetAlign(int64(window))

	data, err := s.client.MultiRangeWithOptions(from, to, options, filter...)
	if err != nil {
		return nil, nil, err
	}

	dat := map[string]map[int64]float64{}
	labels := map[string]map[string]string{}

	for _, datum := range data {
		trace := map[int64]float64{}

		for _, point := range datum.DataPoints {
			trace[point.Timestamp] = point.Value
		}
		labels[datum.Name] = datum.Labels
		dat[datum.Name] = trace
	}

	return dat, labels, nil
}

func (s *Store) Summary(key string, start int64, stop int64, window int, mode string) (map[int64]float64, error) {

	options := redistimeseries.DefaultRangeOptions

	options.SetAggregation(redistimeseries.AggregationType(mode), window)

	data, err := s.client.RangeWithOptions(key, start, stop, options)
	if err != nil {
		return nil, err
	}

	dat := map[int64]float64{}
	for _, datum := range data {
		dat[datum.Timestamp] = datum.Value
	}

	return dat, nil
}

func (s *Store) Push(key string, labels map[string]string, value float64) error {

	r, haveit := s.client.Info(key)
	if haveit != nil {
		op := redistimeseries.DefaultCreateOptions
		op.Labels = labels
		err := s.client.CreateKeyWithOptions(key, op)
		return err
	} else {
		for k := range labels {
			_, ok := r.Labels[k]
			if !ok {
				op := redistimeseries.DefaultCreateOptions
				op.Labels = labels
				_ = s.client.AlterKeyWithOptions(key, op)
			}
		}
	}

	_, err := s.client.AddAutoTs(key, value)
	if err != nil {
		return err
	}

	return nil
}
