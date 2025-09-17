package utils

import (
	"bytes"
	"chat_app/paseto"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Auth tests
func TestCheckAuth(t *testing.T) {
	// Set up test environment
	os.Setenv("PASETO_SYMMETRIC_KEY", "test-symmetric-key-for-testing-only")
	defer os.Unsetenv("PASETO_SYMMETRIC_KEY")

	maker, err := paseto.New()
	if err != nil {
		t.Fatalf("Failed to create paseto maker: %v", err)
	}

	userID := primitive.NewObjectID()
	username := "testuser"
	duration := 24 * time.Hour

	// Create a valid token
	token, err := maker.CreateToken(userID, username, duration)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	tests := []struct {
		name        string
		setupReq    func() *http.Request
		expectError bool
		errorType   string
	}{
		{
			name: "Valid auth cookie",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.AddCookie(&http.Cookie{
					Name:  "auth_cookie",
					Value: token,
				})
				return req
			},
			expectError: false,
		},
		{
			name: "No auth cookie",
			setupReq: func() *http.Request {
				return httptest.NewRequest("GET", "/test", nil)
			},
			expectError: true,
			errorType:   "noCookieToken",
		},
		{
			name: "Invalid auth cookie",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.AddCookie(&http.Cookie{
					Name:  "auth_cookie",
					Value: "invalid-token",
				})
				return req
			},
			expectError: true,
			errorType:   "cookieNotValid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupReq()
			payload, errResp := CheckAuth(req, maker)

			if tt.expectError {
				if errResp == nil {
					t.Error("Expected error response, got nil")
					return
				}
				if errResp.Type != tt.errorType {
					t.Errorf("Expected error type '%s', got '%s'", tt.errorType, errResp.Type)
				}
			} else {
				if errResp != nil {
					t.Errorf("Unexpected error: %v", errResp)
					return
				}
				if payload == nil {
					t.Error("Expected payload, got nil")
					return
				}
				if payload.UserId != userID {
					t.Errorf("Expected user ID %v, got %v", userID, payload.UserId)
				}
				if payload.Username != username {
					t.Errorf("Expected username %s, got %s", username, payload.Username)
				}
			}
		})
	}
}

// JSON tests
func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		data     any
		expected string
	}{
		{
			name:     "Simple string",
			status:   http.StatusOK,
			data:     "test",
			expected: `"test"`,
		},
		{
			name:   "Map data",
			status: http.StatusOK,
			data: map[string]string{
				"key": "value",
			},
			expected: `{
	"key": "value"
}`,
		},
		{
			name:   "Slice data",
			status: http.StatusOK,
			data:   []string{"item1", "item2"},
			expected: `[
	"item1",
	"item2"
]`,
		},
		{
			name:   "Struct data",
			status: http.StatusCreated,
			data: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{
				ID:   1,
				Name: "test",
			},
			expected: `{
	"id": 1,
	"name": "test"
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteJSON(w, tt.status, tt.data)

			if w.Code != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}

			body := w.Body.String()
			if body != tt.expected {
				t.Errorf("Expected body '%s', got '%s'", tt.expected, body)
			}
		})
	}
}

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		maxBytes    int64
		target      any
		expectError bool
	}{
		{
			name:        "Valid JSON object",
			jsonData:    `{"name": "test", "age": 25}`,
			maxBytes:    1024,
			target:      &map[string]any{},
			expectError: false,
		},
		{
			name:        "Valid JSON array",
			jsonData:    `["item1", "item2"]`,
			maxBytes:    1024,
			target:      &[]string{},
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			jsonData:    `{"name": "test", "age": 25,}`,
			maxBytes:    1024,
			target:      &map[string]any{},
			expectError: true,
		},
		{
			name:        "Empty JSON",
			jsonData:    "",
			maxBytes:    1024,
			target:      &map[string]any{},
			expectError: true,
		},
		{
			name:        "Exceeds max bytes",
			jsonData:    `{"data": "very long string that exceeds the limit"}`,
			maxBytes:    10,
			target:      &map[string]any{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := io.NopCloser(bytes.NewBufferString(tt.jsonData))
			err := ParseJSON(body, tt.maxBytes, tt.target)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// Error tests
func TestWriteError(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		errType   any
		errDetail any
		expected  ErrorResponse
	}{
		{
			name:      "String error type",
			status:    http.StatusBadRequest,
			errType:   "validation_error",
			errDetail: "Invalid input",
			expected: ErrorResponse{
				Type:   "validation_error",
				Detail: "Invalid input",
			},
		},
		{
			name:      "Error interface",
			status:    http.StatusInternalServerError,
			errType:   errors.New("server_error"),
			errDetail: "Additional details",
			expected: ErrorResponse{
				Type:   "server_error",
				Detail: "Additional details",
			},
		},
		{
			name:      "Other type",
			status:    http.StatusNotFound,
			errType:   404,
			errDetail: "Resource not found",
			expected: ErrorResponse{
				Type:   "unexepected error type :<nil>",
				Detail: "Resource not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteError(w, tt.status, tt.errType, tt.errDetail)

			if w.Code != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
			}

			var response ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Type != tt.expected.Type {
				t.Errorf("Expected error type '%s', got '%s'", tt.expected.Type, response.Type)
			}

			if response.Detail != tt.expected.Detail {
				t.Errorf("Expected error detail '%v', got '%v'", tt.expected.Detail, response.Detail)
			}
		})
	}
}

// Hash tests
func TestHash(t *testing.T) {
	plainText := []byte("password123")

	hash, err := Hash(plainText)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if len(hash) == 0 {
		t.Error("Expected non-empty hash")
	}

	// Verify the hash can be used to verify the original password
	if !VerifyHash(string(hash), string(plainText)) {
		t.Error("Hash verification should succeed for the same password")
	}

	// Verify wrong password fails
	if VerifyHash(string(hash), "wrongpassword") {
		t.Error("Hash verification should fail for wrong password")
	}
}

// Benchmark tests
func BenchmarkHash(b *testing.B) {
	plainText := []byte("password123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Hash(plainText)
		if err != nil {
			b.Fatalf("Failed to hash password: %v", err)
		}
	}
}

func BenchmarkWriteJSON(b *testing.B) {
	data := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		WriteJSON(w, http.StatusOK, data)
	}
}

func BenchmarkParseJSON(b *testing.B) {
	jsonData := `{"name": "test", "age": 25, "email": "test@example.com"}`
	target := &map[string]any{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body := io.NopCloser(bytes.NewBufferString(jsonData))
		err := ParseJSON(body, 1024, target)
		if err != nil {
			b.Fatalf("Failed to parse JSON: %v", err)
		}
	}
}
