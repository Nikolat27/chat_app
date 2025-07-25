package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type WebSocket struct {
	ChatConnections map[string]map[string]*websocket.Conn
	ConnMu          *sync.Mutex
}

type WsConnection struct {
	Conn *websocket.Conn
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Init() *WebSocket {
	return &WebSocket{
		ChatConnections: make(map[string]map[string]*websocket.Conn),
		ConnMu:          &sync.Mutex{},
	}
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*WsConnection, error) {
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

func (wsConn *WsConnection) Add(chatId, userId string, wsInstance *WebSocket) {
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

func (wsConn *WsConnection) HandleInputs(chatId, userId string, wsInstance *WebSocket) {
	defer func() {
		wsConn.Delete(chatId, userId, wsInstance)
		wsConn.Conn.Close()
	}()

	for {
		messageType, payload, err := wsConn.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		connections := wsInstance.ChatConnections[chatId]

		for _, conn := range connections {
			conn.WriteMessage(messageType, payload)
		}
		
	}
}
