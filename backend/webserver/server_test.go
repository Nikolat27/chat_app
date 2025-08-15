package webserver

import (
	"chat_app/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// Create a mock handler
	mockHandler := &handlers.Handler{}

	tests := []struct {
		name     string
		port     string
		handler  *handlers.Handler
		expected string
	}{
		{"Valid port", "8080", mockHandler, ":8080"},
		{"Port with leading colon", ":8080", mockHandler, "::8080"},
		{"Empty port", "", mockHandler, ":"},
		{"Large port number", "65535", mockHandler, ":65535"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := New(tt.port, tt.handler)

			if server == nil {
				t.Fatal("Expected server instance, got nil")
			}

			if server.Port != tt.port {
				t.Errorf("Expected port %s, got %s", tt.port, server.Port)
			}

			if server.Server == nil {
				t.Fatal("Expected http.Server instance, got nil")
			}

			if server.Server.Addr != tt.expected {
				t.Errorf("Expected server address %s, got %s", tt.expected, server.Server.Addr)
			}

			if server.Server.Handler == nil {
				t.Fatal("Expected handler, got nil")
			}
		})
	}
}

func TestNewWithNilHandler(t *testing.T) {
	server := New("8080", nil)

	if server == nil {
		t.Fatal("Expected server instance, got nil")
	}

	if server.Server == nil {
		t.Fatal("Expected http.Server instance, got nil")
	}

	// Server should still be created even with nil handler
	if server.Server.Handler == nil {
		t.Fatal("Expected handler, got nil")
	}
}

func TestServerConfiguration(t *testing.T) {
	mockHandler := &handlers.Handler{}
	server := New("8080", mockHandler)

	// Test server configuration
	if server.Server.ReadTimeout != 0 {
		t.Errorf("Expected default read timeout 0, got %v", server.Server.ReadTimeout)
	}

	if server.Server.WriteTimeout != 0 {
		t.Errorf("Expected default write timeout 0, got %v", server.Server.WriteTimeout)
	}

	if server.Server.IdleTimeout != 0 {
		t.Errorf("Expected default idle timeout 0, got %v", server.Server.IdleTimeout)
	}

	if server.Server.MaxHeaderBytes != 0 {
		t.Errorf("Expected default max header bytes 0, got %d", server.Server.MaxHeaderBytes)
	}
}

func TestServerWithCustomConfiguration(t *testing.T) {
	mockHandler := &handlers.Handler{}
	server := New("8080", mockHandler)

	// Set custom timeouts
	server.Server.ReadTimeout = 30 * time.Second
	server.Server.WriteTimeout = 30 * time.Second
	server.Server.IdleTimeout = 60 * time.Second
	server.Server.MaxHeaderBytes = 1 << 20 // 1MB

	// Verify custom configuration
	if server.Server.ReadTimeout != 30*time.Second {
		t.Errorf("Expected read timeout 30s, got %v", server.Server.ReadTimeout)
	}

	if server.Server.WriteTimeout != 30*time.Second {
		t.Errorf("Expected write timeout 30s, got %v", server.Server.WriteTimeout)
	}

	if server.Server.IdleTimeout != 60*time.Second {
		t.Errorf("Expected idle timeout 60s, got %v", server.Server.IdleTimeout)
	}

	if server.Server.MaxHeaderBytes != 1<<20 {
		t.Errorf("Expected max header bytes 1MB, got %d", server.Server.MaxHeaderBytes)
	}
}

func TestServerClose(t *testing.T) {
	mockHandler := &handlers.Handler{}
	server := New("8080", mockHandler)

	// Test close on unstarted server
	err := server.Close()
	if err != nil {
		t.Errorf("Expected no error when closing unstarted server, got %v", err)
	}
}

func TestServerRunHttps(t *testing.T) {
	mockHandler := &handlers.Handler{}
	server := New("8080", mockHandler)

	// Test HTTPS with non-existent certificate files
	err := server.RunHttps("nonexistent.crt", "nonexistent.key")
	if err == nil {
		t.Error("Expected error when certificate files don't exist, got nil")
	}
}

func TestServerIntegration(t *testing.T) {
	// Create a simple test handler
	testHandler := &handlers.Handler{}
	server := New("0", testHandler) // Use port 0 for automatic port assignment

	// Start server in a goroutine
	go func() {
		err := server.Run()
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test that server can be closed
	err := server.Close()
	if err != nil {
		t.Errorf("Failed to close server: %v", err)
	}
}

func TestServerWithRouter(t *testing.T) {
	mockHandler := &handlers.Handler{}
	server := New("8080", mockHandler)

	// Test that the server has a router
	if server.Server.Handler == nil {
		t.Fatal("Expected router handler, got nil")
	}

	// Test that the router can handle requests
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// This should not panic
	server.Server.Handler.ServeHTTP(w, req)
}

func TestServerPortValidation(t *testing.T) {
	mockHandler := &handlers.Handler{}

	tests := []struct {
		name        string
		port        string
		expectError bool
	}{
		{"Valid port", "8080", false},
		{"Valid port with colon", ":8080", false},
		{"Zero port", "0", false},
		{"Large port", "65535", false},
		{"Invalid port (negative)", "-1", false},
		{"Invalid port (too large)", "65536", false},
		{"Non-numeric port", "abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := New(tt.port, mockHandler)

			if server == nil {
				t.Fatal("Expected server instance, got nil")
			}

			// Server should be created regardless of port validity
			// The actual port validation happens when the server starts
			if server.Server == nil {
				t.Fatal("Expected http.Server instance, got nil")
			}
		})
	}
}

func TestServerHandlerIntegration(t *testing.T) {
	// Create a mock handler with some basic functionality
	mockHandler := &handlers.Handler{}
	server := New("8080", mockHandler)

	// Test that the server can be configured with different handlers
	if server.Server.Handler == nil {
		t.Fatal("Expected handler, got nil")
	}

	// Test that the handler can be accessed
	handler := server.Server.Handler
	if handler == nil {
		t.Fatal("Handler should not be nil")
	}
}

// Benchmark tests
func BenchmarkNew(b *testing.B) {
	mockHandler := &handlers.Handler{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server := New("8080", mockHandler)
		if server == nil {
			b.Fatal("Expected server instance, got nil")
		}
	}
}

func BenchmarkServerClose(b *testing.B) {
	mockHandler := &handlers.Handler{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server := New("8080", mockHandler)
		err := server.Close()
		if err != nil {
			b.Fatalf("Failed to close server: %v", err)
		}
	}
}

func BenchmarkServerConfiguration(b *testing.B) {
	mockHandler := &handlers.Handler{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		server := New("8080", mockHandler)

		// Set custom configuration
		server.Server.ReadTimeout = 30 * time.Second
		server.Server.WriteTimeout = 30 * time.Second
		server.Server.IdleTimeout = 60 * time.Second
		server.Server.MaxHeaderBytes = 1 << 20

		// Verify configuration
		_ = server.Server.ReadTimeout
		_ = server.Server.WriteTimeout
		_ = server.Server.IdleTimeout
		_ = server.Server.MaxHeaderBytes
	}
}
