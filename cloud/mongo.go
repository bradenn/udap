package cloud

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type Cloud struct {
	database *mongo.Database
}

func New() (Cloud, error) {
	// The credentials are retrieved from the OS environment
	dbUser := os.Getenv("dbUser")
	dbPass := os.Getenv("dbPass")
	// Host and port are also obtained from the environment
	dbHost := os.Getenv("dbHost")
	// The name of the database is again retrieved from the environment
	dbName := os.Getenv("dbName")
	// Now we organize the variables into the connection string
	urlString := fmt.Sprintf("mongodb://%s:%s@%s:27017/%s", dbUser, dbPass, dbHost, dbName)
	// If the client cannot connect within 10 seconds, cancel the attempt
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Apply the connection string to the new client
	clientOptions := options.Client().ApplyURI(urlString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return Cloud{}, err
	}
	database := client.Database(dbName)
	cloud := Cloud{database: database}

	return cloud, nil
}

type Persistent interface {
	Collection() string
}

func (c *Cloud) Insert(persistent Persistent) (string, error) {
	one, err := c.From(persistent.Collection()).InsertOne(context.Background(), persistent)
	if err != nil {
		return "", err
	}
	id := one.InsertedID.(primitive.ObjectID)
	return id.Hex(), nil
}

func (c *Cloud) Update(persistent Persistent) (string, error) {
	// one := c.From(persistent.Collection()).FindOneAndUpdate(context.Background(), bson.persistent)
	// if err != nil {
	// 	return "", err
	// }
	// id := one.InsertedID.(primitive.ObjectID)
	// return id.Hex(), nil
	return "", nil
}

func (c *Cloud) Delete(persistent Persistent) error {
	return nil
}

func (c *Cloud) From(collection string) *mongo.Collection {
	return c.database.Collection(collection)
}
