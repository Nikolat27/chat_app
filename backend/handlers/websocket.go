package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketManager -> Rooms connections
type WebSocketManager struct {
	ChatConnections  map[string]map[string]*websocket.Conn // chatId -> userId -> ws Conn
	GroupConnections map[string]map[string]*websocket.Conn // groupId -> userId -> ws Conn
	UserConnections  map[string]*websocket.Conn            // userId -> ws Conn (ensures 1 connection per user)
	ConnMutex        sync.RWMutex                          // Read/Write mutex for better concurrency
}

// WsConnection -> Websocket connection itself for users
type WsConnection struct {
	Conn *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 3 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		allowedOrigins := getAllowedOrigins()
		currentOrigin := r.Header.Get("Origin")

		return slices.Contains(allowedOrigins, currentOrigin)
	},
}

// WebsocketInit -> Constructor
func WebsocketInit() *WebSocketManager {
	return &WebSocketManager{
		ChatConnections:  make(map[string]map[string]*websocket.Conn),
		GroupConnections: make(map[string]map[string]*websocket.Conn),
		UserConnections:  make(map[string]*websocket.Conn),
		ConnMutex:        sync.RWMutex{},
	}
}

func WebsocketUpgrade(w http.ResponseWriter, r *http.Request) (*WsConnection, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return &WsConnection{
		Conn: conn,
	}, nil
}

// Close -> Close the connection
func (wsConn *WsConnection) Close() {
	if err := wsConn.Conn.Close(); err != nil {
		slog.Error("closing ws conn", "error", err)
	}
}

// GetGoroutineCount -> Get current number of goroutines (for monitoring)
func (ws *WebSocketManager) GetGoroutineCount() int {
	// This is a simple way to get goroutine count
	// In production, you might want to use runtime.NumGoroutine()
	return 0 // Placeholder - you can implement actual goroutine counting if needed
}

// Delete -> Delete user`s connection from either chatId or groupId. Also deletes the room if it`s empty
func (ws *WebSocketManager) Delete(roomId, userId string) {
	ws.ConnMutex.Lock()
	defer ws.ConnMutex.Unlock()

	// Clean up user connection tracker
	if conn, exists := ws.UserConnections[userId]; exists {
		slog.Info("user disconnected", "user_id", userId, "room_id", roomId)
		conn.Close()
		delete(ws.UserConnections, userId)

		// Remove user from all existing rooms (chat or group)
		ws.removeUserFromAllRooms(userId)
	}

	// Try to delete from chat connections first
	if connections, ok := ws.ChatConnections[roomId]; ok {
		delete(connections, userId)
		if len(connections) == 0 {
			delete(ws.ChatConnections, roomId)
			slog.Info("empty chat room deleted", "room_id", roomId)
		}
		return
	}

	// If not found in chat connections, try group connections
	if connections, ok := ws.GroupConnections[roomId]; ok {
		delete(connections, userId)
		if len(connections) == 0 {
			delete(ws.GroupConnections, roomId)
			slog.Info("empty group room deleted", "room_id", roomId)
		}
		return
	}
}

// ChatMessage -> Both Regular and Secret Chats
type ChatMessage struct {
	SenderId       string `json:"sender_id"`
	ReceiverId     string `json:"receiver_id"`
	Content        string `json:"content"`         // content is only for text messages
	ContentAddress string `json:"content_address"` // content address is only for images
	ContentType    string `json:"content_type"`    // either an image or text
}

// AddChat -> Adds the user's websocket connection to the chat room
// Ensures user can only be in one room at a time (either chat or group)
// Maximum 2 users per chat (including the new user)
func (wsConn *WsConnection) AddChat(chatId, userId string, ws *WebSocketManager) {
	ws.ConnMutex.Lock()
	defer ws.ConnMutex.Unlock()

	// Check if user already has a connection and close it (from any room)
	if existingConn, exists := ws.UserConnections[userId]; exists {
		slog.Info("closing existing connection for user", "user_id", userId)
		existingConn.Close()

		// Remove user from all existing rooms (chat or group)
		ws.removeUserFromAllRooms(userId)
	}

	// Check chat connection limit (max 2 users)
	connections, ok := ws.ChatConnections[chatId]
	if !ok {
		connections = make(map[string]*websocket.Conn)
		ws.ChatConnections[chatId] = connections
	} else if len(connections) >= 2 {
		slog.Warn("chat connection limit reached", "chat_id", chatId, "current_users", len(connections))
		return
	}

	// Add new connection
	ws.UserConnections[userId] = wsConn.Conn
	connections[userId] = wsConn.Conn

	slog.Info("user connected to chat", "user_id", userId, "chat_id", chatId, "total_users", len(connections))
}

// HandleChatIncomingMsgs -> for both regular chat and secret chat
func (wsConn *WsConnection) HandleChatIncomingMsgs(chatId, senderId, receiverId string, isSecret bool,
	wsInstance *WebSocketManager, handler *Handler) error {

	defer func() {
		wsInstance.Delete(chatId, senderId)
		if err := wsConn.Conn.Close(); err != nil {
			slog.Error("closing ws conn", "error", err)
		}
		slog.Info("chat websocket handler ended", "chat_id", chatId, "user_id", senderId)
	}()

	slog.Info("chat websocket handler started", "chat_id", chatId, "user_id", senderId, "is_secret", isSecret)

	for {
		_, payload, err := wsConn.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("unexpected websocket close", "error", err, "chat_id", chatId, "user_id", senderId)
			}
			return fmt.Errorf("failed to ws read message: %w", err)
		}

		// Wrap the payload with sender info
		var input ChatMessage

		if err := json.Unmarshal(payload, &input); err != nil {
			slog.Error("failed to unmarshal chat message", "error", err, "chat_id", chatId, "user_id", senderId)
			return fmt.Errorf("failed to UnMarshal ws message: %w", err)
		}

		// Get connections with read lock for better performance
		chatConnections := wsInstance.GetChatConnectionsReadOnly(chatId)
		// to prevent panics...
		if chatConnections == nil {
			slog.Warn("no connections found for chat", "chat_id", chatId)
			continue
		}

		// Store message to DB in background
		go func() {
			if err := handler.storeChatMsgToDB(chatId, senderId, receiverId, input.ContentType,
				input.ContentAddress, input.Content, isSecret); err != nil {
				slog.Error("failed to store chat message to DB", "error", err, "chat_id", chatId, "sender_id", senderId)
			}
		}()

		// Use the optimized broadcast method
		if err := wsInstance.BroadcastToRoom(chatId, senderId, websocket.TextMessage, payload); err != nil {
			slog.Error("failed to broadcast chat message", "error", err, "chat_id", chatId, "sender_id", senderId)
			return fmt.Errorf("failed to broadcast message: %w", err)
		}
	}
}

// GroupMessage -> Both Regular and Secret Groups
type GroupMessage struct {
	SenderId       string `json:"sender_id"`
	Content        string `json:"content"`
	ContentAddress string `json:"content_address"`
	ContentType    string `json:"content_type"` // either an image or text
}

// HandleGroupIncomingMsgs -> Handles both Regular and Secret groups
func (wsConn *WsConnection) HandleGroupIncomingMsgs(groupId, senderId string, isSecret bool, wsInstance *WebSocketManager,
	handler *Handler) error {

	defer func() {
		wsInstance.Delete(groupId, senderId)
		if err := wsConn.Conn.Close(); err != nil {
			slog.Error("closing ws conn", "error", err)
		}
		slog.Info("group websocket handler ended", "group_id", groupId, "user_id", senderId)
	}()

	slog.Info("group websocket handler started", "group_id", groupId, "user_id", senderId, "is_secret", isSecret)

	for {
		_, payload, err := wsConn.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("unexpected websocket close", "error", err, "group_id", groupId, "user_id", senderId)
			}
			return fmt.Errorf("failed to read message: %w", err)
		}

		var input GroupMessage

		if err := json.Unmarshal(payload, &input); err != nil {
			slog.Error("failed to unmarshal group message", "error", err, "group_id", groupId, "user_id", senderId)
			return fmt.Errorf("failed to UnMarshal message: %w", err)
		}

		// Store message to DB in background
		go func() {
			if err := handler.storeGroupMsgToDB(groupId, senderId, input.ContentType, input.ContentAddress,
				input.Content, isSecret); err != nil {
				slog.Error("failed to store group message to DB", "error", err, "group_id", groupId, "sender_id", senderId)
			}
		}()

		// Use the optimized broadcast method
		if err := wsInstance.BroadcastToRoom(groupId, senderId, websocket.TextMessage, payload); err != nil {
			slog.Error("failed to broadcast group message", "error", err, "group_id", groupId, "sender_id", senderId)
			return fmt.Errorf("failed to broadcast message: %w", err)
		}
	}
}

// AddGroup -> Adds the user's websocket connection to the group room
// Ensures user can only be in one room at a time (either chat or group)
// Maximum 100 users per group (including the new user)
func (wsConn *WsConnection) AddGroup(groupId, userId string, ws *WebSocketManager) {
	ws.ConnMutex.Lock()
	defer ws.ConnMutex.Unlock()

	// Check if user already has a connection and close it (from any room)
	if existingConn, exists := ws.UserConnections[userId]; exists {
		slog.Info("closing existing connection for user", "user_id", userId)
		existingConn.Close()

		// Remove user from all existing rooms (chat or group)
		ws.removeUserFromAllRooms(userId)
	}

	// Check group connection limit (max 100 users)
	conns, ok := ws.GroupConnections[groupId]
	if !ok {
		conns = make(map[string]*websocket.Conn)
		ws.GroupConnections[groupId] = conns
	} else if len(conns) >= 100 {
		slog.Warn("group connection limit reached", "group_id", groupId, "current_users", len(conns))
		return
	}

	// Add new connection
	ws.UserConnections[userId] = wsConn.Conn
	conns[userId] = wsConn.Conn

	slog.Info("user connected to group", "user_id", userId, "group_id", groupId, "total_users", len(conns))
}

// GetChatConnections -> Safely get chat connections with mutex protection
func (ws *WebSocketManager) GetChatConnections(chatId string) map[string]*websocket.Conn {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	if connections, ok := ws.ChatConnections[chatId]; ok {
		// Return a copy to avoid race conditions
		result := make(map[string]*websocket.Conn)
		for k, v := range connections {
			result[k] = v
		}
		return result
	}
	return nil
}

// GetGroupConnections -> Safely get group connections with mutex protection
func (ws *WebSocketManager) GetGroupConnections(groupId string) map[string]*websocket.Conn {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	if connections, ok := ws.GroupConnections[groupId]; ok {
		// Return a copy to avoid race conditions
		result := make(map[string]*websocket.Conn)
		for k, v := range connections {
			result[k] = v
		}
		return result
	}
	return nil
}

// GetChatConnectionsReadOnly -> Get chat connections with read lock (for broadcasting)
func (ws *WebSocketManager) GetChatConnectionsReadOnly(chatId string) map[string]*websocket.Conn {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	if connections, ok := ws.ChatConnections[chatId]; ok {
		// Return a copy to avoid race conditions during iteration
		result := make(map[string]*websocket.Conn, len(connections))
		for k, v := range connections {
			result[k] = v
		}
		return result
	}
	return nil
}

// GetGroupConnectionsReadOnly -> Get group connections with read lock (for broadcasting)
func (ws *WebSocketManager) GetGroupConnectionsReadOnly(groupId string) map[string]*websocket.Conn {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	if connections, ok := ws.GroupConnections[groupId]; ok {
		// Return a copy to avoid race conditions during iteration
		result := make(map[string]*websocket.Conn, len(connections))
		for k, v := range connections {
			result[k] = v
		}
		return result
	}
	return nil
}

// IsUserConnectedToChat -> Check if user is already connected to a chat
func (ws *WebSocketManager) IsUserConnectedToChat(chatId, userId string) bool {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	if connections, ok := ws.ChatConnections[chatId]; ok {
		_, exists := connections[userId]
		return exists
	}
	return false
}

// IsUserConnectedToGroup -> Check if user is already connected to a group
func (ws *WebSocketManager) IsUserConnectedToGroup(groupId, userId string) bool {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	if connections, ok := ws.GroupConnections[groupId]; ok {
		_, exists := connections[userId]
		return exists
	}
	return false
}

// IsUserConnected -> Check if user has any active connection
func (ws *WebSocketManager) IsUserConnected(userId string) bool {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	_, exists := ws.UserConnections[userId]
	return exists
}

// GetUserConnection -> Get user's current connection (if any)
func (ws *WebSocketManager) GetUserConnection(userId string) (*websocket.Conn, bool) {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	conn, exists := ws.UserConnections[userId]
	return conn, exists
}

// ForceDisconnectUser -> Force disconnect a user from all connections
func (ws *WebSocketManager) ForceDisconnectUser(userId string) bool {
	ws.ConnMutex.Lock()
	defer ws.ConnMutex.Unlock()

	if conn, exists := ws.UserConnections[userId]; exists {
		slog.Info("force disconnecting user", "user_id", userId)
		conn.Close()
		delete(ws.UserConnections, userId)

		// Remove from all rooms
		for roomId, connections := range ws.ChatConnections {
			if _, userExists := connections[userId]; userExists {
				delete(connections, userId)
				if len(connections) == 0 {
					delete(ws.ChatConnections, roomId)
				}
			}
		}
		for roomId, connections := range ws.GroupConnections {
			if _, userExists := connections[userId]; userExists {
				delete(connections, userId)
				if len(connections) == 0 {
					delete(ws.GroupConnections, roomId)
				}
			}
		}
		return true
	}
	return false
}

// GetConnectionStats -> Get statistics about current connections
func (ws *WebSocketManager) GetConnectionStats() map[string]any {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	stats := make(map[string]any)

	// Count total users connected
	stats["total_users"] = len(ws.UserConnections)

	// Count chat rooms
	totalChatConnections := 0
	for _, connections := range ws.ChatConnections {
		totalChatConnections += len(connections)
	}
	stats["chat_rooms"] = len(ws.ChatConnections)
	stats["chat_connections"] = totalChatConnections

	// Count group rooms
	totalGroupConnections := 0
	for _, connections := range ws.GroupConnections {
		totalGroupConnections += len(connections)
	}
	stats["group_rooms"] = len(ws.GroupConnections)
	stats["group_connections"] = totalGroupConnections

	// Add limits info
	stats["chat_limit"] = 2
	stats["group_limit"] = 100

	return stats
}

// CheckRoomLimits -> Check if any rooms are approaching their limits
func (ws *WebSocketManager) CheckRoomLimits() map[string]any {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	warnings := make(map[string]any)

	// Check chat rooms approaching limit
	for chatId, connections := range ws.ChatConnections {
		if len(connections) >= 2 {
			warnings[fmt.Sprintf("chat_%s", chatId)] = map[string]any{
				"type":    "chat",
				"room_id": chatId,
				"current": len(connections),
				"limit":   2,
				"status":  "at_limit",
			}
		}
	}

	// Check group rooms approaching limit
	for groupId, connections := range ws.GroupConnections {
		if len(connections) >= 100 {
			warnings[fmt.Sprintf("group_%s", groupId)] = map[string]any{
				"type":    "group",
				"room_id": groupId,
				"current": len(connections),
				"limit":   100,
				"status":  "at_limit",
			}
		}
	}

	return warnings
}

// RoomExists -> Quick check if a room exists (chat or group)
func (ws *WebSocketManager) RoomExists(roomId string) (bool, string) {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	if _, exists := ws.ChatConnections[roomId]; exists {
		return true, "chat"
	}
	if _, exists := ws.GroupConnections[roomId]; exists {
		return true, "group"
	}
	return false, ""
}

// GetRoomType -> Get the type of room (chat or group)
func (ws *WebSocketManager) GetRoomType(roomId string) string {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	if _, exists := ws.ChatConnections[roomId]; exists {
		return "chat"
	}
	if _, exists := ws.GroupConnections[roomId]; exists {
		return "group"
	}
	return ""
}

// GetUserRoomInfo -> Get information about which room a user is connected to
func (ws *WebSocketManager) GetUserRoomInfo(userId string) (string, string) {
	ws.ConnMutex.RLock()
	defer ws.ConnMutex.RUnlock()

	// Check chat rooms
	for chatId, connections := range ws.ChatConnections {
		if _, exists := connections[userId]; exists {
			return chatId, "chat"
		}
	}

	// Check group rooms
	for groupId, connections := range ws.GroupConnections {
		if _, exists := connections[userId]; exists {
			return groupId, "group"
		}
	}

	return "", ""
}

// BroadcastToRoom -> Broadcast message to all users in a room except sender
func (ws *WebSocketManager) BroadcastToRoom(roomId, senderId string, messageType int, payload []byte) error {
	var connections map[string]*websocket.Conn
	var roomType string

	// Get connections based on room type
	if ws.GetRoomType(roomId) == "chat" {
		connections = ws.GetChatConnectionsReadOnly(roomId)
		roomType = "chat"
	} else {
		connections = ws.GetGroupConnectionsReadOnly(roomId)
		roomType = "group"
	}

	if connections == nil {
		return fmt.Errorf("room %s (%s) not found or empty", roomId, roomType)
	}

	var errors []string
	successCount := 0

	for userId, conn := range connections {
		if userId == senderId {
			continue
		}

		// Check if connection is still valid
		if conn == nil {
			errors = append(errors, fmt.Sprintf("nil connection for user %s", userId))
			continue
		}

		if err := conn.WriteMessage(messageType, payload); err != nil {
			errors = append(errors, fmt.Sprintf("failed to send message to %s: %v", userId, err))
		} else {
			successCount++
		}
	}

	if len(errors) > 0 {
		slog.Warn("broadcast errors", "room_id", roomId, "room_type", roomType, "errors", errors)
	}

	slog.Debug("broadcast completed", "room_id", roomId, "room_type", roomType,
		"success_count", successCount, "error_count", len(errors))

	if len(errors) > 0 {
		return fmt.Errorf("broadcast errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// removeUserFromAllRooms -> Helper method to remove user from all rooms
func (ws *WebSocketManager) removeUserFromAllRooms(userId string) {
	// Remove from chat rooms
	for roomId, connections := range ws.ChatConnections {
		if _, userExists := connections[userId]; userExists {
			delete(connections, userId)
			if len(connections) == 0 {
				delete(ws.ChatConnections, roomId)
			}
		}
	}

	// Remove from group rooms
	for roomId, connections := range ws.GroupConnections {
		if _, userExists := connections[userId]; userExists {
			delete(connections, userId)
			if len(connections) == 0 {
				delete(ws.GroupConnections, roomId)
			}
		}
	}
}

func getAllowedOrigins() []string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		return []string{"http://localhost:5000"}
	}

	return strings.Split(origins, ",")
}
