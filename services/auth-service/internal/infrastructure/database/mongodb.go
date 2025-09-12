package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoConfig struct {
	URI            string
	Database       string
	ConnectTimeout time.Duration
	PingTimeout    time.Duration
	MaxPoolSize    uint64
}

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
	config   *MongoConfig
}

func NewMongoConfig(uri, database string) *MongoConfig {
	return &MongoConfig{
		URI:            uri,
		Database:       database,
		ConnectTimeout: 10 * time.Second,
		PingTimeout:    5 * time.Second,
		MaxPoolSize:    100,
	}
}

func Connect(config *MongoConfig) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(config.URI).
		SetMaxPoolSize(config.MaxPoolSize).
		SetConnectTimeout(config.ConnectTimeout)

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), config.PingTimeout)
	defer pingCancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("error pinging MongoDB: %v", err)
	}

	database := client.Database(config.Database)

	log.Printf("Connected to MongoDB at %s", config.URI)

	return &MongoDB{
		Client:   client,
		Database: database,
		config:   config,
	}, nil
}

func (m *MongoDB) Close(ctx context.Context) error {
	if m.Client == nil {
		return nil
	}
	if err := m.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %v", err)
	}
	log.Printf("Disconnected from MongoDB at %s", m.config.URI)
	return nil
}

func (m *MongoDB) Ping(ctx context.Context) error {
	if m.Client == nil {
		return fmt.Errorf("MongoDB client'ı mevcut değil")
	}

	return m.Client.Ping(ctx, nil)
}

func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

func (m *MongoDB) IsConnected(ctx context.Context) bool {
	if m.Client == nil {
		return false
	}

	return m.Ping(ctx) == nil
}

// GetDatabaseName database adını döndürür
func (m *MongoDB) GetDatabaseName() string {
	if m.config == nil {
		return ""
	}
	return m.config.Database
}
