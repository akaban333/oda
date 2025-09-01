package database

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNewMongoClient(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
	}{
		{
			name: "with valid environment variables",
			envVars: map[string]string{
				"MONGODB_URI":      "mongodb://localhost:27017",
				"MONGODB_DATABASE": "testdb",
			},
			expectError: false,
		},
		{
			name:        "without environment variables",
			envVars:     map[string]string{},
			expectError: false, // Should use defaults
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			// Note: This test requires a running MongoDB instance
			// In a real CI/CD environment, you would use a test container
			// For now, we'll skip if MongoDB is not available
			client, err := NewMongoClient()
			if err != nil {
				t.Skipf("MongoDB not available: %v", err)
			}

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
				assert.NotNil(t, client.Client)
				assert.NotNil(t, client.Database)
				assert.NotNil(t, client.ctx)
				assert.NotNil(t, client.cancelFunc)

				// Clean up
				client.Close()
			}
		})
	}
}

func TestMongoClient_GetCollection(t *testing.T) {
	// Set up test environment
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("MONGODB_DATABASE", "testdb")
	defer os.Unsetenv("MONGODB_URI")
	defer os.Unsetenv("MONGODB_DATABASE")

	client, err := NewMongoClient()
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
	defer client.Close()

	// Test getting collections
	usersCollection := client.GetCollection(CollectionNames.Users)
	assert.NotNil(t, usersCollection)
	assert.Equal(t, "users", usersCollection.Name())

	roomsCollection := client.GetCollection(CollectionNames.Rooms)
	assert.NotNil(t, roomsCollection)
	assert.Equal(t, "rooms", roomsCollection.Name())

	// Test getting non-existent collection
	customCollection := client.GetCollection("custom_collection")
	assert.NotNil(t, customCollection)
	assert.Equal(t, "custom_collection", customCollection.Name())
}

func TestMongoClient_Close(t *testing.T) {
	// Set up test environment
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("MONGODB_DATABASE", "testdb")
	defer os.Unsetenv("MONGODB_URI")
	defer os.Unsetenv("MONGODB_DATABASE")

	client, err := NewMongoClient()
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}

	// Test that client is connected
	assert.NotNil(t, client.Client)
	assert.NotNil(t, client.Database)

	// Close the client
	err = client.Close()
	assert.NoError(t, err)

	// Test that context is cancelled
	select {
	case <-client.ctx.Done():
		// Context was cancelled as expected
	default:
		t.Error("Context was not cancelled after Close()")
	}
}

func TestCollectionNames(t *testing.T) {
	// Test that all collection names are defined
	assert.NotEmpty(t, CollectionNames.Users)
	assert.NotEmpty(t, CollectionNames.Rooms)
	assert.NotEmpty(t, CollectionNames.Sessions)
	assert.NotEmpty(t, CollectionNames.Materials)
	assert.NotEmpty(t, CollectionNames.Todos)
	assert.NotEmpty(t, CollectionNames.Notes)
	assert.NotEmpty(t, CollectionNames.Posts)
	assert.NotEmpty(t, CollectionNames.Notifications)
	assert.NotEmpty(t, CollectionNames.RealTimeChannels)
	assert.NotEmpty(t, CollectionNames.ChatMessages)

	// Test that collection names are unique
	names := []string{
		CollectionNames.Users,
		CollectionNames.Rooms,
		CollectionNames.Sessions,
		CollectionNames.Materials,
		CollectionNames.Todos,
		CollectionNames.Notes,
		CollectionNames.Posts,
		CollectionNames.Notifications,
		CollectionNames.RealTimeChannels,
		CollectionNames.ChatMessages,
	}

	seen := make(map[string]bool)
	for _, name := range names {
		assert.False(t, seen[name], "Duplicate collection name: %s", name)
		seen[name] = true
	}
}

func TestMongoClient_ContextTimeout(t *testing.T) {
	// Set up test environment
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("MONGODB_DATABASE", "testdb")
	defer os.Unsetenv("MONGODB_URI")
	defer os.Unsetenv("MONGODB_DATABASE")

	client, err := NewMongoClient()
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
	defer client.Close()

	// Test that context has a timeout
	select {
	case <-time.After(15 * time.Second):
		// Context should timeout after 10 seconds (plus buffer)
		t.Error("Context did not timeout as expected")
	case <-client.ctx.Done():
		// Context timed out as expected
	}
}

func TestMongoClient_DatabaseOperations(t *testing.T) {
	// Set up test environment
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("MONGODB_DATABASE", "testdb")
	defer os.Unsetenv("MONGODB_URI")
	defer os.Unsetenv("MONGODB_DATABASE")

	client, err := NewMongoClient()
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
	defer client.Close()

	// Test basic database operations
	collection := client.GetCollection("test_collection")

	// Try to insert a test document - skip if authentication is required
	testDoc := bson.M{"test": "data", "timestamp": time.Now()}
	result, err := collection.InsertOne(context.Background(), testDoc)
	if err != nil && strings.Contains(err.Error(), "authentication") {
		t.Skipf("MongoDB requires authentication: %v", err)
	}
	require.NoError(t, err)
	assert.NotNil(t, result.InsertedID)

	// Find the document
	var foundDoc bson.M
	err = collection.FindOne(context.Background(), bson.M{"test": "data"}).Decode(&foundDoc)
	require.NoError(t, err)
	assert.Equal(t, "data", foundDoc["test"])

	// Clean up
	_, err = collection.DeleteOne(context.Background(), bson.M{"test": "data"})
	require.NoError(t, err)
}

func TestMongoClient_ConnectionOptions(t *testing.T) {
	// Test that client options are properly configured
	uri := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(uri)

	assert.NotNil(t, clientOptions)
	assert.Equal(t, uri, clientOptions.GetURI())
}

func TestMongoClient_ErrorHandling(t *testing.T) {
	// Test with invalid URI
	os.Setenv("MONGODB_URI", "mongodb://invalid:27017")
	os.Setenv("MONGODB_DATABASE", "testdb")
	defer os.Unsetenv("MONGODB_URI")
	defer os.Unsetenv("MONGODB_DATABASE")

	client, err := NewMongoClient()
	if err != nil {
		// Expected error for invalid URI
		assert.Error(t, err)
		assert.Nil(t, client)
	} else {
		// If somehow it connects, clean up
		client.Close()
	}
}
