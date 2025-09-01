package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoClient holds the MongoDB client connection
type MongoClient struct {
	Client     *mongo.Client
	Database   *mongo.Database
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewMongoClient creates a new MongoDB client
func NewMongoClient() (*MongoClient, error) {
	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DATABASE")

	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	if dbName == "" {
		dbName = "studyplatform"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancel()
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		cancel()
		return nil, err
	}

	log.Println("Connected to MongoDB")
	return &MongoClient{
		Client:     client,
		Database:   client.Database(dbName),
		ctx:        ctx,
		cancelFunc: cancel,
	}, nil
}

// Close closes the MongoDB connection
func (m *MongoClient) Close() error {
	if m.cancelFunc != nil {
		m.cancelFunc()
	}

	if m.Client != nil {
		err := m.Client.Disconnect(context.Background())
		if err != nil {
			return err
		}
		log.Println("Disconnected from MongoDB")
	}
	return nil
}

// GetCollection returns a collection from the database
func (m *MongoClient) GetCollection(collectionName string) *mongo.Collection {
	return m.Database.Collection(collectionName)
}

// CollectionNames returns the names of all collections in the database
var CollectionNames = struct {
	Users            string
	Rooms            string
	Sessions         string
	Materials        string
	Todos            string
	Notes            string
	Posts            string
	Notifications    string
	RealTimeChannels string
	ChatMessages     string
}{
	Users:            "users",
	Rooms:            "rooms",
	Sessions:         "sessions",
	Materials:        "materials",
	Todos:            "todos",
	Notes:            "notes",
	Posts:            "posts",
	Notifications:    "notifications",
	RealTimeChannels: "realtime_channels",
	ChatMessages:     "chat_messages",
}
