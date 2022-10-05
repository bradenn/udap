// Copyright (c) 2022 Braden Nicholson

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"udap/internal/core/domain"
	"udap/internal/plugin"
)

var Module NOAA

type NOAA struct {
	plugin.Module
	geomagneticEntityId string
}

func init() {
	config := plugin.Config{
		Name:        "noaa",
		Type:        "module",
		Description: "Retrieve Data from NOAA",
		Version:     "0.0.1",
		Author:      "Braden Nicholson",
	}
	Module.Config = config
}

const Geomagnetic = "https://nomads.ncep.noaa.gov/pub/data/nccf/com/swmf/prod/swmf.20220930/IMF.dat"

func (n *NOAA) Setup() (plugin.Config, error) {
	err := n.UpdateInterval(time.Minute * 60)
	if err != nil {
		return plugin.Config{}, err
	}
	return Module.Config, nil
}

type Gauss [][]float64

func (n *NOAA) pull() error {
	client := http.Client{}
	now := time.Now()
	url := fmt.Sprintf("https://nomads.ncep.noaa.gov/pub/data/nccf/com/swmf/prod/swmf.%d%02d%02d/IMF.dat", now.Year(),
		now.Month(), now.Day())
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return err
	}
	defer client.CloseIdleConnections()

	data := buf.String()

	lines := strings.Split(data, "\n")
	gs := Gauss{}
	for i, line := range lines {
		if i < 8 {
			continue
		}
		ln := strings.Split(line, " ")
		var fa []float64
		for _, s := range ln {
			float, err2 := strconv.ParseFloat(s, 64)
			if err2 != nil {
				continue
			}
			fa = append(fa, float)
		}
		gs = append(gs, fa)
	}

	marshal, err := json.Marshal(gs)
	if err != nil {
		return err
	}

	err = n.Attributes.Update(n.geomagneticEntityId, "gauss", string(marshal), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (n *NOAA) Update() error {
	if n.Ready() {
		err := n.pull()
		if err != nil {
			n.Err(err)
		}
	}
	return nil
}

func (n *NOAA) Run() error {

	entity := domain.Entity{
		Name:   "geomagnetic",
		Type:   "media",
		Module: "noaa",
	}

	err := n.Entities.Register(&entity)
	if err != nil {
		return err
	}

	n.geomagneticEntityId = entity.Id

	geomagnetic := domain.Attribute{
		Key:     "gauss",
		Value:   "{}",
		Request: "{}",
		Type:    "media",
		Channel: make(chan domain.Attribute),
		Order:   0,
		Entity:  n.geomagneticEntityId,
	}

	err = n.Attributes.Register(&geomagnetic)
	if err != nil {
		return err
	}

	return nil
}
