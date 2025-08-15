package models

import (
	"context"
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testDB *mongo.Database
var userModel *UserModel

func setupTestDB(t testing.TB) {
	// Set up test environment variables
	os.Setenv("DATABASE_NAME", "test_database")

	// Connect to test database
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("MongoDB connection failed (skipping tests): %v", err)
		return
	}

	testDB = client.Database("test_database")
	userModel = NewUserModel(testDB)
}

func cleanupTestDB(t testing.TB) {
	if testDB != nil {
		// Clean up test data
		err := testDB.Collection("users").Drop(context.Background())
		if err != nil {
			t.Logf("Failed to drop test collection: %v", err)
		}

		// Disconnect from database
		err = testDB.Client().Disconnect(context.Background())
		if err != nil {
			t.Logf("Failed to disconnect from database: %v", err)
		}
	}
}

func TestNewUserModel(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	if userModel == nil {
		t.Fatal("Expected UserModel instance, got nil")
	}

	if userModel.collection == nil {
		t.Fatal("Expected collection, got nil")
	}
}

func TestUserModelCreate(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	if userModel == nil {
		t.Skip("UserModel not available")
	}

	tests := []struct {
		name           string
		username       string
		hashedPassword string
		salt           string
		expectError    bool
	}{
		{"Valid user", "testuser", "hashedpassword123", "salt123", false},
		{"Empty username", "", "hashedpassword123", "salt123", false}, // MongoDB allows empty strings
		{"Empty hashed password", "testuser", "", "salt123", false},   // MongoDB allows empty strings
		{"Empty salt", "testuser", "hashedpassword123", "", false},    // MongoDB allows empty strings
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use unique usernames to avoid duplicate key errors
			uniqueUsername := tt.username
			if tt.username == "testuser" {
				uniqueUsername = fmt.Sprintf("testuser%d", i)
			}

			userID, err := userModel.Create(uniqueUsername, tt.hashedPassword, tt.salt)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if userID == primitive.NilObjectID {
					t.Error("Expected valid ObjectID, got nil")
				}
			}
		})
	}
}

func TestUserModelCreateDuplicateUsername(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	if userModel == nil {
		t.Skip("UserModel not available")
	}

	// Create first user
	userID1, err := userModel.Create("duplicateuser", "password1", "salt1")
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	// Try to create second user with same username
	userID2, err := userModel.Create("duplicateuser", "password2", "salt2")
	if err == nil {
		t.Error("Expected error for duplicate username, got nil")
		// Clean up the second user if it was created
		if userID2 != primitive.NilObjectID {
			userModel.Delete(bson.M{"_id": userID2})
		}
	}

	// Verify first user still exists
	user, err := userModel.Get(bson.M{"_id": userID1}, bson.M{})
	if err != nil {
		t.Errorf("Failed to retrieve first user: %v", err)
	}
	if user == nil {
		t.Error("First user should still exist")
	}
}

func TestUserModelGet(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	if userModel == nil {
		t.Skip("UserModel not available")
	}

	// Create a test user
	userID, err := userModel.Create("testuser", "hashedpassword", "salt")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name        string
		filter      bson.M
		projection  bson.M
		expectError bool
		expectNil   bool
	}{
		{"Get by ID", bson.M{"_id": userID}, bson.M{}, false, false},
		{"Get by username", bson.M{"username": "testuser"}, bson.M{}, false, false},
		{"Get with projection", bson.M{"_id": userID}, bson.M{"username": 1}, false, false}, // Only inclusion, no exclusion
		{"Get non-existent user", bson.M{"_id": primitive.NewObjectID()}, bson.M{}, true, true},
		{"Get with invalid filter", bson.M{"invalid_field": "value"}, bson.M{}, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := userModel.Get(tt.filter, tt.projection)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if tt.expectNil {
				if user != nil {
					t.Error("Expected nil user, got user")
				}
			} else {
				if user == nil {
					t.Error("Expected user, got nil")
				} else {
					// Verify user data based on projection
					if tt.projection["username"] == 1 {
						// Only username is included in projection
						if user.Username != "testuser" {
							t.Errorf("Expected username 'testuser', got '%s'", user.Username)
						}
						// Other fields should be empty due to projection
						if user.HashedPassword != "" {
							t.Errorf("Expected empty hashed password due to projection, got '%s'", user.HashedPassword)
						}
						if user.Salt != "" {
							t.Errorf("Expected empty salt due to projection, got '%s'", user.Salt)
						}
					} else {
						// Full projection - verify all fields
						if user.Username != "testuser" {
							t.Errorf("Expected username 'testuser', got '%s'", user.Username)
						}
						if user.HashedPassword != "hashedpassword" {
							t.Errorf("Expected hashed password 'hashedpassword', got '%s'", user.HashedPassword)
						}
						if user.Salt != "salt" {
							t.Errorf("Expected salt 'salt', got '%s'", user.Salt)
						}
					}
					if user.Id != userID {
						t.Errorf("Expected user ID %v, got %v", userID, user.Id)
					}
				}
			}
		})
	}
}

func TestUserModelUpdate(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	if userModel == nil {
		t.Skip("UserModel not available")
	}

	// Create a test user
	userID, err := userModel.Create("testuser", "oldpassword", "oldsalt")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name        string
		filter      bson.M
		updates     bson.M
		expectError bool
	}{
		{"Update username", bson.M{"_id": userID}, bson.M{"username": "newusername"}, false},
		{"Update password", bson.M{"_id": userID}, bson.M{"hashed_password": "newpassword"}, false},
		{"Update multiple fields", bson.M{"_id": userID}, bson.M{"username": "updateduser", "salt": "newsalt"}, false},
		{"Update non-existent user", bson.M{"_id": primitive.NewObjectID()}, bson.M{"username": "newusername"}, false},
		{"Update with invalid filter", bson.M{"invalid_field": "value"}, bson.M{"username": "newusername"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := userModel.Update(tt.filter, tt.updates)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Error("Expected update result, got nil")
				}
			}
		})
	}

	// Verify the update worked
	updatedUser, err := userModel.Get(bson.M{"_id": userID}, bson.M{})
	if err != nil {
		t.Fatalf("Failed to retrieve updated user: %v", err)
	}
	if updatedUser.Username != "updateduser" {
		t.Errorf("Expected updated username 'updateduser', got '%s'", updatedUser.Username)
	}
}

func TestUserModelDelete(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	if userModel == nil {
		t.Skip("UserModel not available")
	}

	// Create a test user
	userID, err := userModel.Create("testuser", "password", "salt")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name        string
		filter      bson.M
		expectError bool
	}{
		{"Delete by ID", bson.M{"_id": userID}, false},
		{"Delete non-existent user", bson.M{"_id": primitive.NewObjectID()}, false},
		{"Delete with invalid filter", bson.M{"invalid_field": "value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := userModel.Delete(tt.filter)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == nil {
					t.Error("Expected delete result, got nil")
				}
			}
		})
	}

	// Verify the user was deleted
	_, err = userModel.Get(bson.M{"_id": userID}, bson.M{})
	if err == nil {
		t.Error("Expected error when retrieving deleted user, got nil")
	}
}

func TestUserModelCRUDCycle(t *testing.T) {
	setupTestDB(t)
	defer cleanupTestDB(t)

	if userModel == nil {
		t.Skip("UserModel not available")
	}

	// Create
	userID, err := userModel.Create("cruduser", "password", "salt")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Read
	user, err := userModel.Get(bson.M{"_id": userID}, bson.M{})
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if user.Username != "cruduser" {
		t.Errorf("Expected username 'cruduser', got '%s'", user.Username)
	}

	// Update
	_, err = userModel.Update(bson.M{"_id": userID}, bson.M{"username": "updatedcruduser"})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Read again to verify update
	updatedUser, err := userModel.Get(bson.M{"_id": userID}, bson.M{})
	if err != nil {
		t.Fatalf("Read after update failed: %v", err)
	}
	if updatedUser.Username != "updatedcruduser" {
		t.Errorf("Expected updated username 'updatedcruduser', got '%s'", updatedUser.Username)
	}

	// Delete
	_, err = userModel.Delete(bson.M{"_id": userID})
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	_, err = userModel.Get(bson.M{"_id": userID}, bson.M{})
	if err == nil {
		t.Error("Expected error when retrieving deleted user, got nil")
	}
}

// Benchmark tests
func BenchmarkUserModelCreate(b *testing.B) {
	setupTestDB(b)
	defer cleanupTestDB(b)

	if userModel == nil {
		b.Skip("UserModel not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		username := fmt.Sprintf("benchmarkuser%d", i)
		_, err := userModel.Create(username, "password", "salt")
		if err != nil {
			b.Fatalf("Create failed: %v", err)
		}
	}
}

func BenchmarkUserModelGet(b *testing.B) {
	setupTestDB(b)
	defer cleanupTestDB(b)

	if userModel == nil {
		b.Skip("UserModel not available")
	}

	// Create a test user
	userID, err := userModel.Create("benchmarkuser", "password", "salt")
	if err != nil {
		b.Fatalf("Failed to create test user: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := userModel.Get(bson.M{"_id": userID}, bson.M{})
		if err != nil {
			b.Fatalf("Get failed: %v", err)
		}
	}
}

func BenchmarkUserModelUpdate(b *testing.B) {
	setupTestDB(b)
	defer cleanupTestDB(b)

	if userModel == nil {
		b.Skip("UserModel not available")
	}

	// Create a test user
	userID, err := userModel.Create("benchmarkuser", "password", "salt")
	if err != nil {
		b.Fatalf("Failed to create test user: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := userModel.Update(bson.M{"_id": userID}, bson.M{"username": fmt.Sprintf("updateduser%d", i)})
		if err != nil {
			b.Fatalf("Update failed: %v", err)
		}
	}
}
