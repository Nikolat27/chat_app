package paseto

import (
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNew(t *testing.T) {
	// Set up test environment
	os.Setenv("PASETO_SYMMETRIC_KEY", "test-symmetric-key-for-testing-only")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	maker, err := New()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if maker == nil {
		t.Fatal("Expected maker instance, got nil")
	}

	if maker.paseto == nil {
		t.Fatal("Expected paseto instance, got nil")
	}
}

func TestNewWithoutKey(t *testing.T) {
	// Ensure key is not set
	os.Unsetenv("PASETO_SYMMETRIC_KEY")

	_, err := New()
	if err == nil {
		t.Error("Expected error when PASETO_SYMMETRIC_KEY is not set, got nil")
	}

	expectedError := "PASETO_SYMMETRIC_KEY env variable does not exist"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewWithEmptyKey(t *testing.T) {
	// Set empty key
	os.Setenv("PASETO_SYMMETRIC_KEY", "")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	_, err := New()
	if err == nil {
		t.Error("Expected error when PASETO_SYMMETRIC_KEY is empty, got nil")
	}

	expectedError := "PASETO_SYMMETRIC_KEY env variable does not exist"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestCreateToken(t *testing.T) {
	// Set up test environment
	os.Setenv("PASETO_SYMMETRIC_KEY", "test-symmetric-key-for-testing-only")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	maker, err := New()
	if err != nil {
		t.Fatalf("Failed to create maker: %v", err)
	}

	userID := primitive.NewObjectID()
	username := "testuser"
	duration := 24 * time.Hour

	token, err := maker.CreateToken(userID, username, duration)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token, got empty string")
	}

	// Token should be different each time due to timestamp
	token2, err := maker.CreateToken(userID, username, duration)
	if err != nil {
		t.Fatalf("Failed to create second token: %v", err)
	}

	if token == token2 {
		t.Error("Expected different tokens, got same token")
	}
}

func TestCreateTokenWithDifferentUsers(t *testing.T) {
	// Set up test environment
	os.Setenv("PASETO_SYMMETRIC_KEY", "test-symmetric-key-for-testing-only")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	maker, err := New()
	if err != nil {
		t.Fatalf("Failed to create maker: %v", err)
	}

	userID1 := primitive.NewObjectID()
	userID2 := primitive.NewObjectID()
	username1 := "user1"
	username2 := "user2"
	duration := 24 * time.Hour

	token1, err := maker.CreateToken(userID1, username1, duration)
	if err != nil {
		t.Fatalf("Failed to create token for user1: %v", err)
	}

	token2, err := maker.CreateToken(userID2, username2, duration)
	if err != nil {
		t.Fatalf("Failed to create token for user2: %v", err)
	}

	if token1 == token2 {
		t.Error("Expected different tokens for different users, got same token")
	}
}

func TestCreateTokenWithDifferentDurations(t *testing.T) {
	// Set up test environment
	os.Setenv("PASETO_SYMMETRIC_KEY", "test-symmetric-key-for-testing-only")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	maker, err := New()
	if err != nil {
		t.Fatalf("Failed to create maker: %v", err)
	}

	userID := primitive.NewObjectID()
	username := "testuser"

	tests := []struct {
		name     string
		duration time.Duration
	}{
		{"1 hour", 1 * time.Hour},
		{"24 hours", 24 * time.Hour},
		{"7 days", 7 * 24 * time.Hour},
		{"30 days", 30 * 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := maker.CreateToken(userID, username, tt.duration)
			if err != nil {
				t.Fatalf("Failed to create token: %v", err)
			}

			if token == "" {
				t.Error("Expected non-empty token, got empty string")
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	// Set up test environment
	os.Setenv("PASETO_SYMMETRIC_KEY", "test-symmetric-key-for-testing-only")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	maker, err := New()
	if err != nil {
		t.Fatalf("Failed to create maker: %v", err)
	}

	userID := primitive.NewObjectID()
	username := "testuser"
	duration := 24 * time.Hour

	// Create token
	token, err := maker.CreateToken(userID, username, duration)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Verify token
	payload, err := maker.VerifyToken(token)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}

	if payload == nil {
		t.Fatal("Expected payload, got nil")
	}

	if payload.UserId != userID {
		t.Errorf("Expected user ID %v, got %v", userID, payload.UserId)
	}

	if payload.Username != username {
		t.Errorf("Expected username %s, got %s", username, payload.Username)
	}
}

func TestVerifyTokenWithDifferentKeys(t *testing.T) {
	// Create two makers with different keys
	os.Setenv("PASETO_SYMMETRIC_KEY", "first-key")
	maker1, err := New()
	if err != nil {
		t.Fatalf("Failed to create first maker: %v", err)
	}

	os.Setenv("PASETO_SYMMETRIC_KEY", "second-key")
	maker2, err := New()
	if err != nil {
		t.Fatalf("Failed to create second maker: %v", err)
	}

	// Reset to first key for cleanup
	os.Setenv("PASETO_SYMMETRIC_KEY", "first-key")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	userID := primitive.NewObjectID()
	username := "testuser"
	duration := 24 * time.Hour

	// Create token with first maker
	token, err := maker1.CreateToken(userID, username, duration)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Try to verify with second maker (different key)
	_, err = maker2.VerifyToken(token)
	if err == nil {
		t.Error("Expected error when verifying token with different key, got nil")
	}
}

func TestVerifyInvalidToken(t *testing.T) {
	// Set up test environment
	os.Setenv("PASETO_SYMMETRIC_KEY", "test-symmetric-key-for-testing-only")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	maker, err := New()
	if err != nil {
		t.Fatalf("Failed to create maker: %v", err)
	}

	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{"Empty token", "", true},
		{"Invalid token", "invalid.token.here", true},
		{"Malformed token", "not.a.valid.paseto.token", true},
		{"Random string", "randomstring", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := maker.VerifyToken(tt.token)
			if tt.expectError && err == nil {
				t.Error("Expected error for invalid token, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestTokenExpiration(t *testing.T) {
	// Set up test environment
	os.Setenv("PASETO_SYMMETRIC_KEY", "test-symmetric-key-for-testing-only")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	maker, err := New()
	if err != nil {
		t.Fatalf("Failed to create maker: %v", err)
	}

	userID := primitive.NewObjectID()
	username := "testuser"
	duration := 1 * time.Millisecond // Very short duration

	// Create token
	token, err := maker.CreateToken(userID, username, duration)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	// Try to verify expired token
	_, err = maker.VerifyToken(token)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

func TestCreateAndVerifyMultipleTokens(t *testing.T) {
	// Set up test environment
	os.Setenv("PASETO_SYMMETRIC_KEY", "test-symmetric-key-for-testing-only")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	maker, err := New()
	if err != nil {
		t.Fatalf("Failed to create maker: %v", err)
	}

	userID := primitive.NewObjectID()
	username := "testuser"
	duration := 24 * time.Hour

	// Create multiple tokens
	tokens := make([]string, 5)
	for i := 0; i < 5; i++ {
		token, err := maker.CreateToken(userID, username, duration)
		if err != nil {
			t.Fatalf("Failed to create token %d: %v", i, err)
		}
		tokens[i] = token
	}

	// Verify all tokens
	for i, token := range tokens {
		payload, err := maker.VerifyToken(token)
		if err != nil {
			t.Fatalf("Failed to verify token %d: %v", i, err)
		}

		if payload.UserId != userID {
			t.Errorf("Token %d: Expected user ID %v, got %v", i, userID, payload.UserId)
		}

		if payload.Username != username {
			t.Errorf("Token %d: Expected username %s, got %s", i, username, payload.Username)
		}
	}
}

// Benchmark tests
func BenchmarkCreateToken(b *testing.B) {
	os.Setenv("PASETO_SYMMETRIC_KEY", "test-symmetric-key-for-benchmarking")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	maker, err := New()
	if err != nil {
		b.Fatalf("Failed to create maker: %v", err)
	}

	userID := primitive.NewObjectID()
	username := "benchmarkuser"
	duration := 24 * time.Hour

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := maker.CreateToken(userID, username, duration)
		if err != nil {
			b.Fatalf("Failed to create token: %v", err)
		}
	}
}

func BenchmarkVerifyToken(b *testing.B) {
	os.Setenv("PASETO_SYMMETRIC_KEY", "test-symmetric-key-for-benchmarking")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	maker, err := New()
	if err != nil {
		b.Fatalf("Failed to create maker: %v", err)
	}

	userID := primitive.NewObjectID()
	username := "benchmarkuser"
	duration := 24 * time.Hour

	// Create a token for verification
	token, err := maker.CreateToken(userID, username, duration)
	if err != nil {
		b.Fatalf("Failed to create token: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := maker.VerifyToken(token)
		if err != nil {
			b.Fatalf("Failed to verify token: %v", err)
		}
	}
}

func BenchmarkCreateAndVerifyToken(b *testing.B) {
	os.Setenv("PASETO_SYMMETRIC_KEY", "test-symmetric-key-for-benchmarking")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	maker, err := New()
	if err != nil {
		b.Fatalf("Failed to create maker: %v", err)
	}

	userID := primitive.NewObjectID()
	username := "benchmarkuser"
	duration := 24 * time.Hour

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token, err := maker.CreateToken(userID, username, duration)
		if err != nil {
			b.Fatalf("Failed to create token: %v", err)
		}

		_, err = maker.VerifyToken(token)
		if err != nil {
			b.Fatalf("Failed to verify token: %v", err)
		}
	}
}
