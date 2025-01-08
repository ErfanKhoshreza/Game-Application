package mongo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var (
	instance *DB
	once     sync.Once
)

// New initializes a MongoDB session and returns a singleton instance of DB
func New(uri string, dbName string) (*DB, error) {
	var err error
	once.Do(func() {
		clientOptions := options.Client().ApplyURI(uri)

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Connect to MongoDB
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
			return
		}

		// Ping the MongoDB server
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("MongoDB ping failed: %v", err)
			return
		}

		// Initialize the database instance
		instance = &DB{
			Client:   client,
			Database: client.Database(dbName),
		}
		fmt.Println("Connected to MongoDB")
	})
	return instance, err
}

// GetDBInstance returns the singleton DB instance
func GetDBInstance() *DB {
	if instance == nil {
		log.Fatal("Database not initialized. Call New() first.")
	}
	return instance
}
