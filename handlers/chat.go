package handlers

import (
	"chat_app/utils"
	"chat_app/websocket"
	"net/http"
)

func (handler *Handler) AddChatWebsocket(w http.ResponseWriter, r *http.Request) {
	chatId := r.URL.Query().Get("chat_id")
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "chat id is missing")
		return
	}

	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "user id is missing")
		return
	}

	wsConn, err := websocket.Upgrade(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "upgradingWebsocket", err)
		return
	}
	
	wsConn.Add(chatId, userId, handler.WebSocket)

	go func() {
		wsConn.HandleInputs(chatId, userId, handler.WebSocket)
	}()
}
