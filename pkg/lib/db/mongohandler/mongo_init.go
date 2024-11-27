package mongohandler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBWrapper defines an interface for interacting with the MongoDB database
type MongoDBWrapper interface {
	Collection(name string) *mongo.Collection
	Disconnect()
}

// mongoDBWrapper is the concrete implementation of MongoDBWrapper
type mongoDBWrapper struct {
	client   *mongo.Client
	database *mongo.Database
}

// InitMongoDB sets up the MongoDB connection and returns the wrapper interface
func InitMongoDB() MongoDBWrapper {
	mongoDB, err := newMongoDB()
	if err != nil {
		log.Fatalf("MongoDB initialization failed: %v", err)
	}
	return mongoDB
}

// newMongoDB initializes the MongoDB client and returns an instance of MongoDBWrapper
func newMongoDB() (MongoDBWrapper, error) {
	vi := viper.New()
	vi.AutomaticEnv()

	uri := "mongodb://localhost:27017"
	dbName := "users"

	vi.SetDefault("MONGODB_URI", uri)
	vi.SetDefault("DB_NAME", dbName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database to ensure the connection is valid
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")

	// Return the concrete implementation of MongoDBWrapper
	return &mongoDBWrapper{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

// Collection returns a MongoDB collection from the wrapped database
func (m *mongoDBWrapper) Collection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

// Disconnect gracefully closes the MongoDB connection
func (m *mongoDBWrapper) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect from MongoDB: %v", err)
	} else {
		log.Println("Disconnected from MongoDB")
	}
}
