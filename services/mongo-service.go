package services

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitMongoService(uri string) error {
	fmt.Println(uri)
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	_client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Check if the connection was successful
	err = _client.Ping(context.Background(), nil)
	if err != nil {
		return err
	} else {
		client = _client
		return nil
	}
}
func GetMongoService() *mongo.Client {
	return client
}
