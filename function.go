package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"udap/server"
)

type Function struct {
	Persistent
	Name       string      `json:"name" gorm:"unique"`
	Identifier string      `json:"identifier" gorm:"unique"`
	Module     string      `json:"module"`
	Payload    interface{} `json:"payload" gorm:"type:varchar"`
}

func (f *Function) Create() string {

	col := server.Collection("endpoints")
	one, err := col.InsertOne(context.Background(), f)
	if err != nil {
		return ""
	}
	return one.InsertedID.(primitive.ObjectID).String()
}
