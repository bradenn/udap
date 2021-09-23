package main

import (
	"fmt"
	"udap/cloud"
)

type Demo struct {
	Name string `json:"name"`
}

func (d *Demo) Collection() string {
	return "demos"
}

func RunDemo() {
	c, err := cloud.New()
	if err != nil {
		return
	}
	insert, err := c.Insert(&Demo{})
	if err != nil {
		return
	}

	fmt.Println(insert)
}
