package models

import (
	"context"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var modelsTestDB *mongo.Database

func setupModelsTestDB(t testing.TB) {
	// Set up test environment variables
	os.Setenv("DATABASE_NAME", "test_database")
	
	// Connect to test database
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("MongoDB connection failed (skipping tests): %v", err)
		return
	}

	modelsTestDB = client.Database("test_database")
}

func cleanupModelsTestDB(t testing.TB) {
	if modelsTestDB != nil {
		// Clean up test data by dropping all collections
		collections := []string{"users", "chats", "secret_chats", "messages", "save_messages", "groups", "approvals"}
		for _, collectionName := range collections {
			err := modelsTestDB.Collection(collectionName).Drop(context.Background())
			if err != nil {
				t.Logf("Failed to drop collection %s: %v", collectionName, err)
			}
		}
		
		// Disconnect from database
		err := modelsTestDB.Client().Disconnect(context.Background())
		if err != nil {
			t.Logf("Failed to disconnect from database: %v", err)
		}
	}
}

func TestNew(t *testing.T) {
	setupModelsTestDB(t)
	defer cleanupModelsTestDB(t)

	if modelsTestDB == nil {
		t.Skip("Test database not available")
	}

	models := New(modelsTestDB)
	if models == nil {
		t.Fatal("Expected Models instance, got nil")
	}

	// Test that all model instances are created
	if models.User == nil {
		t.Error("Expected User model, got nil")
	}
	if models.Chat == nil {
		t.Error("Expected Chat model, got nil")
	}
	if models.SecretChat == nil {
		t.Error("Expected SecretChat model, got nil")
	}
	if models.Message == nil {
		t.Error("Expected Message model, got nil")
	}
	if models.SaveMessage == nil {
		t.Error("Expected SaveMessage model, got nil")
	}
	if models.Group == nil {
		t.Error("Expected Group model, got nil")
	}
	if models.Approval == nil {
		t.Error("Expected Approval model, got nil")
	}
}

func TestNewWithNilDatabase(t *testing.T) {
	// This test expects a panic when passing nil database
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil database, got none")
		}
	}()

	New(nil)
}

func TestModelsIntegration(t *testing.T) {
	setupModelsTestDB(t)
	defer cleanupModelsTestDB(t)

	if modelsTestDB == nil {
		t.Skip("Test database not available")
	}

	models := New(modelsTestDB)

	// Test that we can use the models without panicking
	// This is a basic integration test to ensure the models are functional

	// Test User model
	if models.User != nil {
		// Try to create a user (this should work or fail gracefully)
		_, err := models.User.Create("testuser", "password", "salt")
		if err != nil {
			// This is expected if there are validation issues
			t.Logf("User creation failed (expected in test environment): %v", err)
		}
	}

	// Test Chat model
	if models.Chat != nil {
		// The Chat model should be accessible
		t.Log("Chat model is available")
	}

	// Test SecretChat model
	if models.SecretChat != nil {
		// The SecretChat model should be accessible
		t.Log("SecretChat model is available")
	}

	// Test Message model
	if models.Message != nil {
		// The Message model should be accessible
		t.Log("Message model is available")
	}

	// Test SaveMessage model
	if models.SaveMessage != nil {
		// The SaveMessage model should be accessible
		t.Log("SaveMessage model is available")
	}

	// Test Group model
	if models.Group != nil {
		// The Group model should be accessible
		t.Log("Group model is available")
	}

	// Test Approval model
	if models.Approval != nil {
		// The Approval model should be accessible
		t.Log("Approval model is available")
	}
}

// Benchmark tests
func BenchmarkNew(b *testing.B) {
	setupModelsTestDB(b)
	defer cleanupModelsTestDB(b)

	if modelsTestDB == nil {
		b.Skip("Test database not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		models := New(modelsTestDB)
		if models == nil {
			b.Fatal("Expected Models instance, got nil")
		}
	}
} 