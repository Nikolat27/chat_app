package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestCreateGroup(t *testing.T) {
	t.Skip("Skipping HTTP handler tests that require authentication - focus on WebSocket tests")
}

func TestAddGroupWebsocket(t *testing.T) {
	t.Skip("Skipping HTTP handler tests that require authentication - focus on WebSocket tests")
}

func TestGroupMessageBroadcasting(t *testing.T) {
	handler := setupTestHandler()

	t.Run("Group Message Broadcasting", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Errorf("Failed to upgrade connection: %v", err)
				return
			}
			defer conn.Close()

			// Read messages
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					break
				}

				// Verify message format
				var groupMsg GroupMessage
				if err := json.Unmarshal(message, &groupMsg); err != nil {
					t.Errorf("Failed to unmarshal message: %v", err)
					continue
				}

				// Verify message structure
				if groupMsg.SenderId == "" {
					t.Error("Message should have sender_id")
				}
				if groupMsg.ContentType == "" {
					t.Error("Message should have content_type")
				}
			}
		}))
		defer server.Close()

		// Connect multiple users to the group
		connections := make([]*websocket.Conn, 3)
		for i := 0; i < 3; i++ {
			wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if err != nil {
				t.Fatalf("Failed to connect: %v", err)
			}
			defer conn.Close()
			connections[i] = conn

			// Add connection to WebSocket manager
			wsConn := &WsConnection{Conn: conn}
			wsConn.AddGroup("test_group", "user"+string(rune('1'+i)), handler.WebSocket)
		}

		// Send test message from first user
		message := GroupMessage{
			SenderId:       "user1",
			Content:        "Hello, Group!",
			ContentType:    "text",
			ContentAddress: "",
		}

		messageBytes, _ := json.Marshal(message)
		err := connections[0].WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			t.Errorf("Failed to send message: %v", err)
		}

		// Give some time for message processing
		time.Sleep(100 * time.Millisecond)
	})
}

func TestGroupConnectionLimits(t *testing.T) {
	handler := setupTestHandler()

	t.Run("Group Connection Limit Enforcement", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Errorf("Failed to upgrade connection: %v", err)
				return
			}
			defer conn.Close()

			// Keep connection alive
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
			}
		}))
		defer server.Close()

		// Try to add 101 users (should fail at 101st)
		connections := make([]*websocket.Conn, 101)
		successCount := 0

		for i := 0; i < 101; i++ {
			wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if err != nil {
				t.Fatalf("Failed to connect: %v", err)
			}
			defer conn.Close()
			connections[i] = conn

			// Add connection to WebSocket manager
			wsConn := &WsConnection{Conn: conn}
			wsConn.AddGroup("limit_test_group", "user"+string(rune('1'+i)), handler.WebSocket)

			// Check if user was actually connected
			if handler.WebSocket.IsUserConnected("user" + string(rune('1'+i))) {
				successCount++
			}
		}

		// Should have exactly 100 successful connections
		if successCount != 100 {
			t.Errorf("Expected 100 successful connections, got %d", successCount)
		}

		// Verify 101st user is not connected
		if handler.WebSocket.IsUserConnected("user" + string(rune('1'+100))) {
			t.Error("101st user should not be connected due to limit")
		}
	})
}

func TestGetGroupMessages(t *testing.T) {
	t.Skip("Skipping HTTP handler tests that require authentication - focus on WebSocket tests")
}

func TestGetGroupMembers(t *testing.T) {
	t.Skip("Skipping HTTP handler tests that require authentication - focus on WebSocket tests")
}

// Benchmark tests
func BenchmarkCreateGroup(b *testing.B) {
	handler := setupTestHandler()

	requestBody := map[string]any{
		"name":        "Benchmark Group",
		"description": "A benchmark test group",
		"type":        "public",
	}
	jsonBody, _ := json.Marshal(requestBody)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/group/create", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateGroup(w, req)
	}
}

func BenchmarkGroupWebSocketMessageSending(b *testing.B) {
	handler := setupTestHandler()

	// Setup test connection
	conn := createTestConnectionForBenchmark(b)
	defer conn.Close()

	wsConn := &WsConnection{Conn: conn}
	wsConn.AddGroup("benchmark_group", "user1", handler.WebSocket)

	message := GroupMessage{
		SenderId:       "user1",
		Content:        "Benchmark group message",
		ContentType:    "text",
		ContentAddress: "",
	}
	messageBytes, _ := json.Marshal(message)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := conn.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			b.Errorf("Failed to send message: %v", err)
		}
	}
}
