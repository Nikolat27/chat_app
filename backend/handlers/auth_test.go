package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Skip("Skipping Register test since it requires database connection")
}

func TestLogin(t *testing.T) {
	t.Skip("Skipping Login test since it requires database connection")
}

func TestLogout(t *testing.T) {
	handler := setupTestHandler()

	t.Run("Successful Logout", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/auth/logout", nil)
		w := httptest.NewRecorder()

		handler.Logout(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Check if cookie is set to be deleted
		cookies := w.Result().Cookies()
		found := false
		for _, cookie := range cookies {
			if cookie.Name == "auth_cookie" && cookie.MaxAge == -1 {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected auth_cookie to be set for deletion")
		}
	})
}

func TestAuthCheck(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/auth/check", nil)
		w := httptest.NewRecorder()

		handler.AuthCheck(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("With Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/auth/check", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.AuthCheck(w, req)

		// Should still fail because the token is invalid in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

// Benchmark tests
func BenchmarkRegister(b *testing.B) {
	handler := setupTestHandler()

	requestBody := map[string]string{
		"username": "benchmarkuser",
		"password": "benchmarkpassword123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)
	}
}

func BenchmarkLogin(b *testing.B) {
	handler := setupTestHandler()

	requestBody := map[string]string{
		"username": "benchmarkuser",
		"password": "benchmarkpassword123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Login(w, req)
	}
}

func BenchmarkLogout(b *testing.B) {
	handler := setupTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/auth/logout", nil)
		w := httptest.NewRecorder()

		handler.Logout(w, req)
	}
}
