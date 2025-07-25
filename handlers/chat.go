package handlers

import (
	"chat_app/utils"
	"chat_app/websocket"
	"fmt"
	"net/http"
)

func (handler *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "checkAuth", err)
		return
	}

	userObjectId, err := utils.ToObjectId(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "strToObjectId", err)
		return
	}

	var input struct {
		SecondUser string `json:"second_user"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err)
		return
	}

	handler.Models.Chat.Create()
	
	fmt.Println(userObjectId)

}

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
