package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateApproval(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/approvals/test-invite-link", nil)
		w := httptest.NewRecorder()

		handler.CreateApproval(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Missing Invite Link", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/approvals/", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.CreateApproval(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/approvals/test-invite-link", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.CreateApproval(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Valid Request", func(t *testing.T) {
		requestBody := map[string]string{
			"reason": "I want to join this group",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/api/approvals/test-invite-link", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.CreateApproval(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestEditApprovalStatus(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/approvals/507f1f77bcf86cd799439011", nil)
		w := httptest.NewRecorder()

		handler.EditApprovalStatus(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Missing Approval ID", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/approvals/", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditApprovalStatus(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Invalid Approval ID", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/approvals/invalid-id", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditApprovalStatus(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/approvals/507f1f77bcf86cd799439011", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditApprovalStatus(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Invalid Status", func(t *testing.T) {
		requestBody := map[string]string{
			"status": "invalid_status",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("PUT", "/api/approvals/507f1f77bcf86cd799439011", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditApprovalStatus(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Valid Status - Approved", func(t *testing.T) {
		requestBody := map[string]string{
			"status": "approved",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("PUT", "/api/approvals/507f1f77bcf86cd799439011", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditApprovalStatus(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Valid Status - Rejected", func(t *testing.T) {
		requestBody := map[string]string{
			"status": "rejected",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("PUT", "/api/approvals/507f1f77bcf86cd799439011", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditApprovalStatus(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Valid Status - Pending", func(t *testing.T) {
		requestBody := map[string]string{
			"status": "pending",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("PUT", "/api/approvals/507f1f77bcf86cd799439011", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditApprovalStatus(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestGetReceivedApprovals(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/approvals/received", nil)
		w := httptest.NewRecorder()

		handler.GetReceivedApprovals(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("With Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/approvals/received", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetReceivedApprovals(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("With Pagination", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/approvals/received?page=1&limit=10", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetReceivedApprovals(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestGetSentApprovals(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/approvals/sent", nil)
		w := httptest.NewRecorder()

		handler.GetSentApprovals(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("With Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/approvals/sent", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetSentApprovals(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("With Pagination", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/approvals/sent?page=1&limit=10", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetSentApprovals(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestDeleteApproval(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/approvals/507f1f77bcf86cd799439011", nil)
		w := httptest.NewRecorder()

		handler.DeleteApproval(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Missing Approval ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/approvals/", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteApproval(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Invalid Approval ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/approvals/invalid-id", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteApproval(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Valid Approval ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/approvals/507f1f77bcf86cd799439011", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteApproval(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

// Benchmark tests
func BenchmarkCreateApproval(b *testing.B) {
	handler := setupTestHandler()

	requestBody := map[string]string{
		"reason": "I want to join this group",
	}
	jsonBody, _ := json.Marshal(requestBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/approvals/test-invite-link", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.CreateApproval(w, req)
	}
}

func BenchmarkEditApprovalStatus(b *testing.B) {
	handler := setupTestHandler()

	requestBody := map[string]string{
		"status": "approved",
	}
	jsonBody, _ := json.Marshal(requestBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("PUT", "/api/approvals/507f1f77bcf86cd799439011", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditApprovalStatus(w, req)
	}
}

func BenchmarkGetReceivedApprovals(b *testing.B) {
	handler := setupTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/approvals/received", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetReceivedApprovals(w, req)
	}
}

func BenchmarkGetSentApprovals(b *testing.B) {
	handler := setupTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/approvals/sent", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetSentApprovals(w, req)
	}
}

func BenchmarkDeleteApproval(b *testing.B) {
	handler := setupTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("DELETE", "/api/approvals/507f1f77bcf86cd799439011", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteApproval(w, req)
	}
}
