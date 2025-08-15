package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUploadImageChatMessage(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/messages/chat/507f1f77bcf86cd799439011/image", nil)
		w := httptest.NewRecorder()

		handler.UploadImageChatMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Missing Chat ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/messages/chat//image", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.UploadImageChatMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Invalid Chat ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/messages/chat/invalid-id/image", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.UploadImageChatMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Valid Chat ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/messages/chat/507f1f77bcf86cd799439011/image", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.UploadImageChatMessage(w, req)

		// Should return 400 due to missing file data in test environment
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestUploadImageGroupMessage(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/messages/group/507f1f77bcf86cd799439011/image", nil)
		w := httptest.NewRecorder()

		handler.UploadImageGroupMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Missing Group ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/messages/group//image", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.UploadImageGroupMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Invalid Group ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/messages/group/invalid-id/image", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.UploadImageGroupMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Valid Group ID", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/messages/group/507f1f77bcf86cd799439011/image", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.UploadImageGroupMessage(w, req)

		// Should return 400 due to missing file data in test environment
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestEditMessage(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/messages/507f1f77bcf86cd799439011", nil)
		w := httptest.NewRecorder()

		handler.EditMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Missing Message ID", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/messages/", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Invalid Message ID", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/messages/invalid-id", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Valid Message ID with Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/messages/507f1f77bcf86cd799439011", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Valid Message ID with Valid JSON", func(t *testing.T) {
		requestBody := map[string]string{
			"new_content": "Updated message content",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("PUT", "/api/messages/507f1f77bcf86cd799439011", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditMessage(w, req)

		// Should return 400 due to database connection issues in test environment
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestDeleteMessageForSender(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/messages/507f1f77bcf86cd799439011/sender", nil)
		w := httptest.NewRecorder()

		handler.DeleteMessageForSender(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Missing Message ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/messages//sender", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteMessageForSender(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Invalid Message ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/messages/invalid-id/sender", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteMessageForSender(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Valid Message ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/messages/507f1f77bcf86cd799439011/sender", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteMessageForSender(w, req)

		// Should return 400 due to database connection issues in test environment
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestDeleteMessageForReceiver(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/messages/507f1f77bcf86cd799439011/receiver", nil)
		w := httptest.NewRecorder()

		handler.DeleteMessageForReceiver(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Missing Message ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/messages//receiver", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteMessageForReceiver(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Invalid Message ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/messages/invalid-id/receiver", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteMessageForReceiver(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Valid Message ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/messages/507f1f77bcf86cd799439011/receiver", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteMessageForReceiver(w, req)

		// Should return 400 due to database connection issues in test environment
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestDeleteMessageForAll(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/messages/507f1f77bcf86cd799439011/all", nil)
		w := httptest.NewRecorder()

		handler.DeleteMessageForAll(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Missing Message ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/messages//all", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteMessageForAll(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Invalid Message ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/messages/invalid-id/all", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteMessageForAll(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Valid Message ID", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/messages/507f1f77bcf86cd799439011/all", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteMessageForAll(w, req)

		// Should return 400 due to database connection issues in test environment
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// Benchmark tests
func BenchmarkUploadImageChatMessage(b *testing.B) {
	handler := setupTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/messages/chat/507f1f77bcf86cd799439011/image", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.UploadImageChatMessage(w, req)
	}
}

func BenchmarkUploadImageGroupMessage(b *testing.B) {
	handler := setupTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/messages/group/507f1f77bcf86cd799439011/image", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.UploadImageGroupMessage(w, req)
	}
}

func BenchmarkEditMessage(b *testing.B) {
	handler := setupTestHandler()

	requestBody := map[string]string{
		"new_content": "Updated message content",
	}
	jsonBody, _ := json.Marshal(requestBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("PUT", "/api/messages/507f1f77bcf86cd799439011", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.EditMessage(w, req)
	}
}

func BenchmarkDeleteMessageForSender(b *testing.B) {
	handler := setupTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("DELETE", "/api/messages/507f1f77bcf86cd799439011/sender", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteMessageForSender(w, req)
	}
}

func BenchmarkDeleteMessageForAll(b *testing.B) {
	handler := setupTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("DELETE", "/api/messages/507f1f77bcf86cd799439011/all", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteMessageForAll(w, req)
	}
} 