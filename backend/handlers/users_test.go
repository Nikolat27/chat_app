package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUploadAvatar(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/users/avatar", nil)
		w := httptest.NewRecorder()

		handler.UploadAvatar(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("With Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/users/avatar", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.UploadAvatar(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/users", nil)
		w := httptest.NewRecorder()

		handler.DeleteUser(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("With Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/users", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.DeleteUser(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestGetUserChats(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/chats", nil)
		w := httptest.NewRecorder()

		handler.GetUserChats(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("With Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/chats", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetUserChats(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestGetUserSecretChats(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/secret-chats", nil)
		w := httptest.NewRecorder()

		handler.GetUserSecretChats(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("With Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/secret-chats", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetUserSecretChats(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

func TestGetUserGroups(t *testing.T) {
	handler := setupTestHandler()

	t.Run("No Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/groups", nil)
		w := httptest.NewRecorder()

		handler.GetUserGroups(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("With Auth Cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/groups", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetUserGroups(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Secret Groups", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/groups?secret=true", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetUserGroups(w, req)

		// Should return 401 due to invalid auth token in test environment
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

// Test helper functions
func TestGetOtherUserId(t *testing.T) {
	t.Run("Valid Participants", func(t *testing.T) {
		user1 := primitive.NewObjectID()
		user2 := primitive.NewObjectID()
		user3 := primitive.NewObjectID()
		participants := []primitive.ObjectID{user1, user2, user3}

		otherId := getOtherUserId(participants, user1)

		if otherId == user1 {
			t.Error("Expected different user ID, got same ID")
		}
	})

	t.Run("Single Participant", func(t *testing.T) {
		user1 := primitive.NewObjectID()
		participants := []primitive.ObjectID{user1}

		otherId := getOtherUserId(participants, user1)

		if otherId != user1 {
			t.Error("Expected same user ID when only one participant")
		}
	})
}

// Benchmark tests
func BenchmarkSearchUser(b *testing.B) {
	handler := setupTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/users/search?q=testuser", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.SearchUser(w, req)
	}
}

func BenchmarkGetUser(b *testing.B) {
	handler := setupTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/users/507f1f77bcf86cd799439011", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetUser(w, req)
	}
}

func BenchmarkGetUserChats(b *testing.B) {
	handler := setupTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/users/chats", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetUserChats(w, req)
	}
}

func BenchmarkGetUserGroups(b *testing.B) {
	handler := setupTestHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/users/groups", nil)
		req.AddCookie(createMockAuthCookie())
		w := httptest.NewRecorder()

		handler.GetUserGroups(w, req)
	}
}
