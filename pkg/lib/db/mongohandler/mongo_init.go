package mongohandler

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ConnectionTimeoutInSecond time.Duration = 3 * time.Second
)

type HealthFn func(context.Context) error

type config struct {
	MONGODB_URI string
	DB_NAME     string
}

// MongoDBWrapper is the concrete implementation of MongoDBWrapper
type MongoDBWrapper struct {
	Database *mongo.Database
	conf     config
	client   *mongo.Client
}

func New() *MongoDBWrapper {
	vi := viper.New()
	vi.AutomaticEnv()

	vi.SetDefault("MONGODB_URI", "mongodb://127.0.0.1:27017")
	vi.SetDefault("DB_NAME", "users")

	return &MongoDBWrapper{
		conf: config{
			MONGODB_URI: vi.GetString("MONGODB_URI"),
			DB_NAME:     vi.GetString("DB_NAME"),
		},
	}
}

// InitMongoDB sets up the MongoDB connection. If error occurs exit the program.
func (m *MongoDBWrapper) InitMongoDB() {
	err := m.initDB()
	if err != nil {
		slog.Error("MongoDB initialization failed", slog.Any("error", err))
		os.Exit(1)
	}
}

// initDB initializes the MongoDB. Return error.
func (m *MongoDBWrapper) initDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeoutInSecond)
	defer cancel()

	// Initialize MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.conf.MONGODB_URI))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database to ensure the connection is valid
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	m.client = client
	m.Database = client.Database(m.conf.DB_NAME)

	slog.InfoContext(ctx, fmt.Sprintf("Connected to MongoDB with %s", m.conf.MONGODB_URI))
	return nil
}

// Collection returns a MongoDB collection from the wrapped database
func (m *MongoDBWrapper) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

// Disconnect gracefully closes the MongoDB connection
func (m *MongoDBWrapper) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.client.Disconnect(ctx); err != nil {
		slog.InfoContext(ctx, "Failed to disconnect from MongoDB", slog.Any("error", err))
	} else {
		slog.InfoContext(ctx, "Disconnected from MongoDB")
	}
}

func (m *MongoDBWrapper) HealthChecker() HealthFn {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		return m.client.Ping(ctx, nil)
	}
}
