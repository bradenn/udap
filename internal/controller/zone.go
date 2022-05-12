// Copyright (c) 2021 Braden Nicholson

package controller

import (
	"encoding/json"
	"sync"
	"udap/internal/bond"
	"udap/internal/log"
	"udap/internal/models"
	"udap/internal/store"
)

type Zones struct {
	PolyBuffer
	Observable
}

func (z *Zones) Handle(event bond.Msg) (res interface{}, err error) {
	switch event.Operation {
	case "compile":
		return z.Compile()
	case "create":
		return z.create(event)
	case "delete":
		return z.delete(event)
	case "restore":
		return z.restore(event)
	}
	return nil, nil
}

func (z *Zones) delete(msg bond.Msg) (res interface{}, err error) {
	zone := z.Find(msg.Id)

	err = zone.Delete()
	if err != nil {
		return nil, err
	}
	z.Set(msg.Id, zone)

	return nil, err
}

func (z *Zones) restore(msg bond.Msg) (res interface{}, err error) {
	zone := z.Find(msg.Id)

	err = zone.Restore()
	if err != nil {
		return nil, err
	}
	z.Set(msg.Id, zone)

	return nil, err
}

func (z *Zones) create(msg bond.Msg) (res interface{}, err error) {
	zone := models.Zone{}
	err = json.Unmarshal([]byte(msg.Payload), &zone)
	if err != nil {
		return nil, err
	}
	register, err := z.Register(zone)
	if err != nil {
		return nil, err
	}

	return register, err
}

func (z *Zones) Register(zone models.Zone) (res *models.Zone, err error) {
	err = zone.Emplace()
	if err != nil {
		return nil, err
	}
	z.set(zone.Id, &zone)
	z.emit(zone.Id, &zone)
	return nil, nil
}

func (z *Zones) Compile() (res []models.Zone, err error) {
	for _, s := range z.Keys() {
		zone := z.Find(s)
		if zone == nil {
			continue
		}

		res = append(res, *zone)
	}
	return res, nil
}

func (z *Zones) EmitAll() (err error) {

	for _, k := range z.Keys() {
		find := z.Find(k)
		z.emit(k, find)
	}

	return nil
}

func (z *Zones) FetchAll() {
	var zones []models.Zone
	err := store.DB.Model(&models.Zone{}).Preload("Entities").Find(&zones).Error
	if err != nil {
		log.Err(err)
		return
	}
	for _, zone := range zones {
		z.set(zone.Id, &zone)
		z.emit(zone.Id, &zone)
	}
}

func (z *Zones) Find(id string) *models.Zone {
	return z.get(id).(*models.Zone)
}

func LoadZones() (m *Zones) {
	m = &Zones{}
	m.data = sync.Map{}
	m.Run()
	m.FetchAll()
	return m
}

func (z *Zones) Set(id string, zone *models.Zone) {
	z.set(id, zone)
	z.emit(zone.Id, &zone)
}

func (z *Zones) Delete(id string) {
	err := z.remove(id)
	if err != nil {
		return
	}

}
