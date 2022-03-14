// Copyright (c) 2021 Braden Nicholson

package models

import (
	"time"
	"udap/internal/store"
)

type Network struct {
	store.Persistent
	Name   string `json:"name"`
	Dns    string `json:"dns"`
	Router string `json:"index"`
	Lease  string `json:"lease"`
	Mask   string `json:"mask"`
	Range  string `json:"range"`
}

// Emplace gets a module from its path
func (n *Network) Emplace() (err error) {
	n.UpdatedAt = time.Now()
	err = store.DB.Model(&Network{}).Where("name = ?", n.Name).FirstOrCreate(&n).Error
	if err != nil {
		return err
	}
	return nil
}

func (n *Network) FetchAll() []Network {
	var networks []Network
	store.DB.Model(&Network{}).Find(&networks)
	return networks
}
