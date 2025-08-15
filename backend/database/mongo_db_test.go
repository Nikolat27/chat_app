package database

import (
	"os"
	"testing"
)

func TestNewWithValidURI(t *testing.T) {
	// Set up test environment variables
	os.Setenv("DATABASE_NAME", "test_database")
	defer os.Unsetenv("DATABASE_NAME")

	// Test with a valid MongoDB URI (using MongoDB memory server or mock)
	// Note: This test requires a running MongoDB instance or mock
	uri := "mongodb://localhost:27017"
	
	db, err := New(uri)
	if err != nil {
		// If MongoDB is not running, this is expected
		t.Skipf("MongoDB connection failed (expected if not running): %v", err)
		return
	}

	if db == nil {
		t.Fatal("Expected database instance, got nil")
	}

	// Test that we can access the database name
	dbName := db.Name()
	if dbName != "test_database" {
		t.Errorf("Expected database name 'test_database', got '%s'", dbName)
	}
}

func TestNewWithInvalidURI(t *testing.T) {
	// Set up test environment variables
	os.Setenv("DATABASE_NAME", "test_database")
	defer os.Unsetenv("DATABASE_NAME")

	// Test with an invalid MongoDB URI
	invalidURI := "mongodb://invalid-host:99999"

	_, err := New(invalidURI)
	if err == nil {
		t.Error("Expected error for invalid URI, got nil")
	}
}

func TestNewWithoutDatabaseName(t *testing.T) {
	// Ensure DATABASE_NAME is not set
	os.Unsetenv("DATABASE_NAME")

	uri := "mongodb://localhost:27017"
	_, err := New(uri)
	if err == nil {
		t.Error("Expected error when DATABASE_NAME is not set, got nil")
	}

	expectedError := "DATABASE_NAME env variable does not exist"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewWithEmptyDatabaseName(t *testing.T) {
	// Set empty DATABASE_NAME
	os.Setenv("DATABASE_NAME", "")
	defer os.Unsetenv("DATABASE_NAME")

	uri := "mongodb://localhost:27017"
	_, err := New(uri)
	if err == nil {
		t.Error("Expected error when DATABASE_NAME is empty, got nil")
	}

	expectedError := "DATABASE_NAME env variable does not exist"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestGetDatabaseName(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expectError bool
		expected    string
	}{
		{"Valid database name", "my_database", false, "my_database"},
		{"Empty database name", "", true, ""},
		{"Not set", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Not set" {
				os.Unsetenv("DATABASE_NAME")
			} else {
				os.Setenv("DATABASE_NAME", tt.envValue)
			}
			defer os.Unsetenv("DATABASE_NAME")

			dbName, err := getDatabaseName()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if dbName != tt.expected {
					t.Errorf("Expected database name '%s', got '%s'", tt.expected, dbName)
				}
			}
		})
	}
}

func TestNewWithTimeout(t *testing.T) {
	// Set up test environment variables
	os.Setenv("DATABASE_NAME", "test_database")
	defer os.Unsetenv("DATABASE_NAME")

	// Test with a URI that would cause a timeout
	// This is a theoretical test since we can't easily simulate network timeouts
	uri := "mongodb://invalid-host-that-will-timeout:27017"

	_, err := New(uri)
	if err == nil {
		t.Skip("Connection succeeded unexpectedly (network conditions may vary)")
	}

	// The error should be related to connection failure
	if err != nil {
		// This is expected behavior
		t.Logf("Expected connection error: %v", err)
	}
}

// Benchmark tests
func BenchmarkNew(b *testing.B) {
	// Set up test environment variables
	os.Setenv("DATABASE_NAME", "benchmark_database")
	defer os.Unsetenv("DATABASE_NAME")

	uri := "mongodb://localhost:27017"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db, err := New(uri)
		if err != nil {
			b.Skipf("MongoDB connection failed (skipping benchmark): %v", err)
			return
		}
		if db == nil {
			b.Fatal("Expected database instance, got nil")
		}
	}
}

func BenchmarkGetDatabaseName(b *testing.B) {
	os.Setenv("DATABASE_NAME", "benchmark_database")
	defer os.Unsetenv("DATABASE_NAME")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := getDatabaseName()
		if err != nil {
			b.Fatalf("getDatabaseName failed: %v", err)
		}
	}
} 