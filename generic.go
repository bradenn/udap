package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Object map[string]interface{}

type ObjectId primitive.ObjectID

type Persistent struct {
	Id string `json:"id" bson:"_id"`
}

func ParseObjectId(id interface{}) string {
	return id.(primitive.ObjectID).Hex()
}
