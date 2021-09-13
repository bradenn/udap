package server

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var database *mongo.Database

func Connect() {
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
		panic(err)
	}
	database = client.Database(dbName)
}

func Collection(collection string) *mongo.Collection {
	return database.Collection(collection)
}
