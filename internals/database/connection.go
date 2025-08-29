package database

import (	
	"context"	// managing deadlines, cancellation signals
	"fmt"	// formatting strings and errors
	"log"	// logging message
	"time"	// dealing with time and durations

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client holds the MongoDB client connection so other functions can use it
var Client *mongo.Client

// ConnectDB establishes connection to MongoDB
func ConnectDB(uri string) error {
	// create context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// set up mongodb client options
	clientOptions := options.Client().ApplyURI(uri)
	// connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	Client = client
	log.Println("✅ connected to MongoDB successfully")
	return nil
}

// GetCollection returns a collection from a database
func GetCollection(databaseName, collectionName string) *mongo.Collection {
	return Client.Database(databaseName).Collection(collectionName)
}

// DisconnectDB	closes the MongoDB connection
func DisconnectDB() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := Client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		} else {
			log.Println("Disconnected from MongoDB")
		}
	}
}
