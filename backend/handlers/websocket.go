package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
    "slices"
	"github.com/gorilla/websocket"
)

// WebSocketManager -> Rooms connections
type WebSocketManager struct {
	ChatConnections  map[string]map[string]*websocket.Conn // chatId -> userId -> ws Conn
	GroupConnections map[string]map[string]*websocket.Conn // groupId -> userId -> ws Conn
	ConnMutex        sync.Mutex
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

		if slices.Contains(allowedOrigins, currentOrigin) {
			return true
		}

		return false
	},
}

// WebsocketInit -> Constructor
func WebsocketInit() *WebSocketManager {
	return &WebSocketManager{
		ChatConnections:  make(map[string]map[string]*websocket.Conn),
		GroupConnections: make(map[string]map[string]*websocket.Conn),
		ConnMutex:        sync.Mutex{},
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

// Delete -> Delete user`s connection from either chatId or groupId. Also deletes the room if it`s empty
func (ws *WebSocketManager) Delete(roomId, userId string) {
	ws.ConnMutex.Lock()
	defer ws.ConnMutex.Unlock()

	// Try to delete from chat connections first
	if connections, ok := ws.ChatConnections[roomId]; ok {
		delete(connections, userId)
		if len(connections) == 0 {
			delete(ws.ChatConnections, roomId)
		}
		return
	}

	// If not found in chat connections, try group connections
	if connections, ok := ws.GroupConnections[roomId]; ok {
		delete(connections, userId)
		if len(connections) == 0 {
			delete(ws.GroupConnections, roomId)
		}
		return
	}
}

// broadcastMessage -> Sends a message to all users in the room except the sender
func broadcastMessage(senderId string, payload []byte, connections map[string]*websocket.Conn) error {
	var errors []string

	for userId, conn := range connections {
		if userId == senderId {
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			errors = append(errors, fmt.Sprintf("failed to send message to %s: %v", userId, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("broadcast errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// ChatMessage -> Both Regular and Secret Chats
type ChatMessage struct {
	SenderId       string `json:"sender_id"`
	ReceiverId     string `json:"receiver_id"`
	Content        string `json:"content"`         // text
	ContentAddress string `json:"content_address"` // image address
	ContentType    string `json:"content_type"`    // either an image or text
}

// AddChat -> Adds the user's websocket connection to the chat room
func (wsConn *WsConnection) AddChat(chatId, userId string, ws *WebSocketManager) {
	ws.ConnMutex.Lock()
	defer ws.ConnMutex.Unlock()

	connections, ok := ws.ChatConnections[chatId]
	if !ok {
		connections = make(map[string]*websocket.Conn)
		ws.ChatConnections[chatId] = connections
	}

	connections[userId] = wsConn.Conn
}

// HandleChatIncomingMsgs -> for both regular chat and secret chat
func (wsConn *WsConnection) HandleChatIncomingMsgs(chatId, senderId, receiverId string, isSecret bool,
	wsInstance *WebSocketManager, handler *Handler) error {

	defer func() {
		wsInstance.Delete(chatId, senderId)
		if err := wsConn.Conn.Close(); err != nil {
			slog.Error("closing ws conn", "error", err)
		}
	}()

	for {
		_, payload, err := wsConn.Conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed to ws read message: %w", err)
		}

		// Wrap the payload with sender info
		var input ChatMessage

		if err := json.Unmarshal(payload, &input); err != nil {
			return fmt.Errorf("failed to UnMarshal ws message: %w", err)
		}

		chatConnections := wsInstance.GetChatConnections(chatId)
		// to prevent panics...
		if chatConnections == nil {
			continue
		}

		if err := handler.storeChatMsgToDB(chatId, senderId, receiverId, input.ContentType,
			input.ContentAddress, input.Content, isSecret); err != nil {

			return fmt.Errorf("failed to store msg in the DB: %w", err)
		}

		if err := broadcastMessage(senderId, payload, chatConnections); err != nil {
			return err
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
	}()

	for {
		_, payload, err := wsConn.Conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		var input GroupMessage

		if err := json.Unmarshal(payload, &input); err != nil {
			return fmt.Errorf("failed to UnMarshal message: %w", err)
		}

		chatConnections := wsInstance.GetGroupConnections(groupId)

		if err := handler.storeGroupMsgToDB(groupId, senderId, input.ContentType, input.ContentAddress,
			input.Content, isSecret); err != nil {
			return fmt.Errorf("failed to store msg in the DB: %w", err)
		}

		if err := broadcastMessage(senderId, payload, chatConnections); err != nil {
			return fmt.Errorf("failed to broadcast msg in the ws connection: %w", err)
		}
	}
}

// AddGroup -> Adds the user's websocket connection to the group room
func (wsConn *WsConnection) AddGroup(groupId, userId string, ws *WebSocketManager) {
	ws.ConnMutex.Lock()
	defer ws.ConnMutex.Unlock()

	conns, ok := ws.GroupConnections[groupId]
	if !ok {
		conns = make(map[string]*websocket.Conn)
		ws.GroupConnections[groupId] = conns
	}

	conns[userId] = wsConn.Conn
}

// GetChatConnections -> Safely get chat connections with mutex protection
func (ws *WebSocketManager) GetChatConnections(chatId string) map[string]*websocket.Conn {
	ws.ConnMutex.Lock()
	defer ws.ConnMutex.Unlock()

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
	ws.ConnMutex.Lock()
	defer ws.ConnMutex.Unlock()

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

// IsUserConnectedToChat -> Check if user is already connected to a chat
func (ws *WebSocketManager) IsUserConnectedToChat(chatId, userId string) bool {
	ws.ConnMutex.Lock()
	defer ws.ConnMutex.Unlock()

	if connections, ok := ws.ChatConnections[chatId]; ok {
		_, exists := connections[userId]
		return exists
	}
	return false
}

// IsUserConnectedToGroup -> Check if user is already connected to a group
func (ws *WebSocketManager) IsUserConnectedToGroup(groupId, userId string) bool {
	ws.ConnMutex.Lock()
	defer ws.ConnMutex.Unlock()

	if connections, ok := ws.GroupConnections[groupId]; ok {
		_, exists := connections[userId]
		return exists
	}
	return false
}

func getAllowedOrigins() []string {
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		return []string{"http://localhost:5000"}
	}

	return strings.Split(origins, ",")
}
