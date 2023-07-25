package services

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func GetMongoService() *mongo.Client {
	if client != nil {
		return client
	}

	uri := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	_client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the connection was successful
	err = _client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	fmt.Println("Connected to MongoDB!")
	client = _client

	return client
}
