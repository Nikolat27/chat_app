package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type WebSocket struct {
	ChatConnections  map[string]map[string]*websocket.Conn
	GroupConnections map[string]map[string]*websocket.Conn
	ConnMu           *sync.Mutex
}

type WsConnection struct {
	Conn *websocket.Conn
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebsocketInit() *WebSocket {
	return &WebSocket{
		ChatConnections:  make(map[string]map[string]*websocket.Conn),
		GroupConnections: make(map[string]map[string]*websocket.Conn),
		ConnMu:           &sync.Mutex{},
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

func (wsConn *WsConnection) Close() {
	wsConn.Conn.Close()
}

func (wsConn *WsConnection) AddChat(chatId, userId string, wsInstance *WebSocket) {
	wsInstance.ConnMu.Lock()
	defer wsInstance.ConnMu.Unlock()

	if wsInstance.ChatConnections[chatId] == nil {
		wsInstance.ChatConnections[chatId] = make(map[string]*websocket.Conn)
	}

	wsInstance.ChatConnections[chatId][userId] = wsConn.Conn
}

func (wsConn *WsConnection) Delete(chatId, userId string, wsInstance *WebSocket) {
	wsInstance.ConnMu.Lock()
	defer wsInstance.ConnMu.Unlock()

	delete(wsInstance.ChatConnections[chatId], userId)
	if len(wsInstance.ChatConnections[chatId]) == 0 {
		delete(wsInstance.ChatConnections, chatId)
	}
}

// Chats
type ChatMessage struct {
	MessageType int    `json:"message_type"`
	SenderId    string `json:"sender_id"`
	ReceiverId  string `json:"receiver_id"`
	Content     string `json:"content"`
}

// both regular chat and secret chat
func (wsConn *WsConnection) HandleChatIncomingMsgs(chatId, senderId, receiverId string, isSecret bool, wsInstance *WebSocket,
	handler *Handler) error {

	defer func() {
		wsConn.Delete(chatId, senderId, wsInstance)
		wsConn.Conn.Close()
	}()

	for {
		messageType, payload, err := wsConn.Conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed to read message: %s", err)
		}

		// Wrap the payload with sender info
		msg := ChatMessage{
			MessageType: messageType,
			SenderId:    senderId,
			ReceiverId:  receiverId,
			Content:     string(payload),
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %s", err)
		}

		chatConnections := wsInstance.ChatConnections[chatId]

		if err := handler.storeChatMsgToDB(chatId, senderId, receiverId, isSecret, payload); err != nil {
			return fmt.Errorf("failed to store msg in the DB: %s", err)
		}

		for userId, conn := range chatConnections {
			if userId == senderId {
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
				return fmt.Errorf("failed to send message to %s: %v\n", userId, err)
			}
		}
	}
}

// Groups
type GroupMessage struct {
	MessageType int    `json:"message_type"`
	SenderId    string `json:"sender_id"`
	Content     string `json:"content"`
}

func (wsConn *WsConnection) HandleGroupIncomingMsgs(groupId, senderId string, isSecret bool, wsInstance *WebSocket,
	handler *Handler) error {

	defer func() {
		wsConn.Delete(groupId, senderId, wsInstance)
		wsConn.Conn.Close()
	}()

	for {
		messageType, payload, err := wsConn.Conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed to read message: %s", err)
		}

		// Wrap the payload with sender info
		msg := GroupMessage{
			MessageType: messageType,
			SenderId:    senderId,
			Content:     string(payload),
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %s", err)
		}

		chatConnections := wsInstance.GroupConnections[groupId]

		if err := handler.storeGroupMsgToDB(groupId, senderId, isSecret, payload); err != nil {
			return fmt.Errorf("failed to store msg in the DB: %s", err)
		}

		for userId, conn := range chatConnections {
			if userId == senderId {
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
				return fmt.Errorf("failed to send message to %s: %v\n", userId, err)
			}
		}
	}
}

func (wsConn *WsConnection) AddGroup(groupId, userId string, wsInstance *WebSocket) {
	wsInstance.ConnMu.Lock()
	defer wsInstance.ConnMu.Unlock()

	if wsInstance.GroupConnections[groupId] == nil {
		wsInstance.GroupConnections[groupId] = make(map[string]*websocket.Conn)
	}

	wsInstance.GroupConnections[groupId][userId] = wsConn.Conn
}

// Secret Groups
type SecretGroupMessage struct {
	MessageType        int               `json:"message_type"`
	SenderId           string            `json:"sender_id"`
	Content            string            `json:"content"`
	UsersSymmetricKeys map[string]string `json:"users_symmetric_keys"`
}

func (wsConn *WsConnection) HandleSecretGroupIncomingMsgs(groupId, senderId string, wsInstance *WebSocket,
	handler *Handler) error {

	defer func() {
		wsConn.Delete(groupId, senderId, wsInstance)
		wsConn.Conn.Close()
	}()

	for {
		messageType, payload, err := wsConn.Conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed to read message: %s", err)
		}

		var input struct {
			Content            string            `json:"content"`
			UsersSymmetricKeys map[string]string `json:"users_symmetric_keys"`
		}

		if err := json.Unmarshal(payload, &input); err != nil {
			return fmt.Errorf("failed to unMarshal message: %s", err)
		}

		// Wrap the payload with sender info
		msg := SecretGroupMessage{
			MessageType:        messageType,
			SenderId:           senderId,
			Content:            input.Content,
			UsersSymmetricKeys: input.UsersSymmetricKeys,
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %s", err)
		}

		chatConnections := wsInstance.GroupConnections[groupId]

		if err := handler.storeSecretGroupMsgToDB(groupId, senderId, []byte(input.Content), input.UsersSymmetricKeys); err != nil {
			return fmt.Errorf("failed to store msg in the DB: %s", err)
		}

		for userId, conn := range chatConnections {
			if userId == senderId {
				continue
			}

			if err := conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
				return fmt.Errorf("failed to send message to %s: %v\n", userId, err)
			}
		}
	}
}
