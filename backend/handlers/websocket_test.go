package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebSocketManager_ConnectionLimits(t *testing.T) {
	ws := WebsocketInit()

	// Test chat connection limit (max 2 users)
	t.Run("Chat Connection Limit", func(t *testing.T) {
		// Create test connections
		conn1 := createTestConnection(t)
		conn2 := createTestConnection(t)
		conn3 := createTestConnection(t)

		// Add first two users (should succeed)
		wsConn1 := &WsConnection{Conn: conn1}
		wsConn2 := &WsConnection{Conn: conn2}

		wsConn1.AddChat("chat1", "user1", ws)
		wsConn2.AddChat("chat1", "user2", ws)

		// Verify both users are connected
		if !ws.IsUserConnected("user1") {
			t.Error("user1 should be connected")
		}
		if !ws.IsUserConnected("user2") {
			t.Error("user2 should be connected")
		}

		// Try to add third user (should fail due to limit)
		wsConn3 := &WsConnection{Conn: conn3}
		wsConn3.AddChat("chat1", "user3", ws)

		// Verify third user is not connected
		if ws.IsUserConnected("user3") {
			t.Error("user3 should not be connected due to limit")
		}

		// Cleanup
		conn1.Close()
		conn2.Close()
		conn3.Close()
	})

	// Test group connection limit (max 100 users)
	t.Run("Group Connection Limit", func(t *testing.T) {
		// Create test connections for 101 users
		connections := make([]*websocket.Conn, 101)
		for i := 0; i < 101; i++ {
			connections[i] = createTestConnection(t)
		}

		// Add first 100 users (should succeed)
		for i := 0; i < 100; i++ {
			wsConn := &WsConnection{Conn: connections[i]}
			wsConn.AddGroup("group1", fmt.Sprintf("user%d", i), ws)
		}

		// Verify first 100 users are connected
		for i := 0; i < 100; i++ {
			if !ws.IsUserConnected(fmt.Sprintf("user%d", i)) {
				t.Errorf("user%d should be connected", i)
			}
		}

		// Try to add 101st user (should fail due to limit)
		wsConn101 := &WsConnection{Conn: connections[100]}
		wsConn101.AddGroup("group1", "user101", ws)

		// Verify 101st user is not connected
		if ws.IsUserConnected("user101") {
			t.Error("user101 should not be connected due to limit")
		}

		// Cleanup
		for _, conn := range connections {
			conn.Close()
		}
	})
}

func TestWebSocketManager_OneConnectionPerUser(t *testing.T) {
	ws := WebsocketInit()

	t.Run("User Can Only Be In One Room", func(t *testing.T) {
		// Create test connections
		conn1 := createTestConnection(t)
		conn2 := createTestConnection(t)

		// Add user to chat
		wsConn1 := &WsConnection{Conn: conn1}
		wsConn1.AddChat("chat1", "user1", ws)

		// Verify user is in chat
		if !ws.IsUserConnected("user1") {
			t.Error("user1 should be connected to chat")
		}

		// Add same user to group (should disconnect from chat)
		wsConn2 := &WsConnection{Conn: conn2}
		wsConn2.AddGroup("group1", "user1", ws)

		// Verify user is now in group, not chat
		if !ws.IsUserConnected("user1") {
			t.Error("user1 should be connected to group")
		}

		// Verify user is not in chat anymore
		if ws.IsUserConnectedToChat("chat1", "user1") {
			t.Error("user1 should not be in chat anymore")
		}

		// Cleanup
		conn1.Close()
		conn2.Close()
	})
}

func TestWebSocketManager_MessageBroadcasting(t *testing.T) {
	ws := WebsocketInit()

	t.Run("Chat Message Broadcasting", func(t *testing.T) {
		// Create test connections
		conn1 := createTestConnection(t)
		conn2 := createTestConnection(t)

		// Add users to chat
		wsConn1 := &WsConnection{Conn: conn1}
		wsConn2 := &WsConnection{Conn: conn2}

		wsConn1.AddChat("chat1", "user1", ws)
		wsConn2.AddChat("chat1", "user2", ws)

		// Create test message
		message := ChatMessage{
			SenderId:       "user1",
			ReceiverId:     "user2",
			Content:        "Hello, World!",
			ContentType:    "text",
			ContentAddress: "",
		}

		messageBytes, _ := json.Marshal(message)

		// Broadcast message
		err := ws.BroadcastToRoom("chat1", "user1", websocket.TextMessage, messageBytes)
		if err != nil {
			t.Errorf("Failed to broadcast message: %v", err)
		}

		// Verify message was sent (in a real test, you'd read from the connection)
		// For this test, we just verify no error occurred

		// Cleanup
		conn1.Close()
		conn2.Close()
	})
}

func TestWebSocketManager_ConnectionCleanup(t *testing.T) {
	ws := WebsocketInit()

	t.Run("User Disconnection Cleanup", func(t *testing.T) {
		// Create test connection
		conn := createTestConnection(t)

		// Add user to chat
		wsConn := &WsConnection{Conn: conn}
		wsConn.AddChat("chat1", "user1", ws)

		// Verify user is connected
		if !ws.IsUserConnected("user1") {
			t.Error("user1 should be connected")
		}

		// Disconnect user
		ws.Delete("chat1", "user1")

		// Verify user is not connected
		if ws.IsUserConnected("user1") {
			t.Error("user1 should not be connected after disconnect")
		}

		// Verify chat room is empty
		connections := ws.GetChatConnectionsReadOnly("chat1")
		if connections != nil && len(connections) > 0 {
			t.Error("chat room should be empty after user disconnect")
		}

		// Cleanup
		conn.Close()
	})
}

func TestWebSocketManager_Statistics(t *testing.T) {
	ws := WebsocketInit()

	t.Run("Connection Statistics", func(t *testing.T) {
		// Create test connections
		conn1 := createTestConnection(t)
		conn2 := createTestConnection(t)
		conn3 := createTestConnection(t)

		// Add users to different rooms
		wsConn1 := &WsConnection{Conn: conn1}
		wsConn2 := &WsConnection{Conn: conn2}
		wsConn3 := &WsConnection{Conn: conn3}

		wsConn1.AddChat("chat1", "user1", ws)
		wsConn2.AddChat("chat2", "user2", ws)
		wsConn3.AddGroup("group1", "user3", ws)

		// Get statistics
		stats := ws.GetConnectionStats()

		// Verify statistics
		if stats["total_users"] != 3 {
			t.Errorf("Expected 3 total users, got %v", stats["total_users"])
		}
		if stats["chat_rooms"] != 2 {
			t.Errorf("Expected 2 chat rooms, got %v", stats["chat_rooms"])
		}
		if stats["group_rooms"] != 1 {
			t.Errorf("Expected 1 group room, got %v", stats["group_rooms"])
		}
		if stats["chat_connections"] != 2 {
			t.Errorf("Expected 2 chat connections, got %v", stats["chat_connections"])
		}
		if stats["group_connections"] != 1 {
			t.Errorf("Expected 1 group connection, got %v", stats["group_connections"])
		}

		// Cleanup
		conn1.Close()
		conn2.Close()
		conn3.Close()
	})
}

func TestWebSocketManager_RoomLimits(t *testing.T) {
	ws := WebsocketInit()

	t.Run("Room Limits Monitoring", func(t *testing.T) {
		// Create test connections for a full chat
		conn1 := createTestConnection(t)
		conn2 := createTestConnection(t)

		// Add users to chat (should reach limit)
		wsConn1 := &WsConnection{Conn: conn1}
		wsConn2 := &WsConnection{Conn: conn2}

		wsConn1.AddChat("chat1", "user1", ws)
		wsConn2.AddChat("chat1", "user2", ws)

		// Check room limits
		warnings := ws.CheckRoomLimits()

		// Verify chat is at limit
		found := false
		for _, warning := range warnings {
			if warningMap, ok := warning.(map[string]any); ok {
				if warningMap["room_id"] == "chat1" && warningMap["status"] == "at_limit" {
					found = true
					break
				}
			}
		}

		if !found {
			t.Error("Chat should be reported as at limit")
		}

		// Cleanup
		conn1.Close()
		conn2.Close()
	})
}

// Helper function to create a test WebSocket connection
func createTestConnection(t *testing.T) *websocket.Conn {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upgrade to WebSocket
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
		// Keep connection open for test duration
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}))
	defer server.Close()

	// Connect to the test server
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to test server: %v", err)
	}

	return conn
}

// Helper function to create a test WebSocket connection for benchmarks
func createTestConnectionForBenchmark(b *testing.B) *websocket.Conn {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upgrade to WebSocket
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			b.Errorf("Failed to upgrade connection: %v", err)
			return
		}
		// Keep connection open for test duration
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}))
	defer server.Close()

	// Connect to the test server
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		b.Fatalf("Failed to connect to test server: %v", err)
	}

	return conn
}

// Benchmark tests for performance
func BenchmarkWebSocketManager_AddChat(b *testing.B) {
	ws := WebsocketInit()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn := createTestConnectionForBenchmark(b)
		wsConn := &WsConnection{Conn: conn}
		wsConn.AddChat(fmt.Sprintf("chat%d", i), fmt.Sprintf("user%d", i), ws)
		conn.Close()
	}
}

func BenchmarkWebSocketManager_Broadcast(b *testing.B) {
	ws := WebsocketInit()

	// Setup test connections
	connections := make([]*websocket.Conn, 10)
	for i := 0; i < 10; i++ {
		connections[i] = createTestConnectionForBenchmark(b)
		wsConn := &WsConnection{Conn: connections[i]}
		wsConn.AddGroup("group1", fmt.Sprintf("user%d", i), ws)
	}
	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()

	message := []byte("test message")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ws.BroadcastToRoom("group1", "user1", websocket.TextMessage, message)
	}
}
