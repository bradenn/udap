// Copyright (c) 2021 Braden Nicholson

package client

import "udap/internal/models"

type Request struct {
}

type Client struct {
	request chan Request
	update  chan Request
}

func (c *Client) UpdateEntity(entity *models.Entity) {

}

func (c *Client) AddEntity(entity *models.Entity) {

}

func (c *Client) RemoveEntity(entity *models.Entity) {

}
