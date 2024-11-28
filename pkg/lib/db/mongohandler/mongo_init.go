package mongohandler

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ConnectionTimeoutInSeconds time.Duration = 3 * time.Second
)

// MongoDBWrapper defines an interface for interacting with the MongoDB database
type MongoDBWrapper interface {
	Collection(name string) Collection
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

	vi.SetDefault("MONGODB_URI", "mongodb://127.0.0.1:27017")
	vi.SetDefault("DB_NAME", "users")
	uri := vi.GetString("MONGODB_URI")
	dbName := vi.GetString("DB_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeoutInSeconds)
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

	slog.InfoContext(ctx, fmt.Sprintf("Connected to MongoDB with %s", uri))

	// Return the concrete implementation of MongoDBWrapper
	return &mongoDBWrapper{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

// Collection returns a MongoDB collection from the wrapped database
func (m *mongoDBWrapper) Collection(name string) Collection {
	return &collection{
		m.database.Collection(name),
	}
}

// Disconnect gracefully closes the MongoDB connection
func (m *mongoDBWrapper) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		slog.InfoContext(ctx, "Failed to disconnect from MongoDB", slog.Any("error", err))
	} else {
		slog.InfoContext(ctx, "Disconnected from MongoDB")
	}
}
