package handlers

import (
	"chat_app/utils"
	"chat_app/websocket"
	"errors"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		TargetUser string `json:"target_user"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err)
		return
	}

	targetUserObjectId, err := utils.ToObjectId(input.TargetUser)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "strToObjectId", err)
		return
	}

	participants := []primitive.ObjectID{userObjectId, targetUserObjectId}

	filter := bson.M{
		"participants": bson.M{
			"$all": participants,
		},
	}

	projection := bson.M{
		"_id": 1,
	}

	// checking chat duplication
	if _, err := handler.Models.Chat.Get(filter, projection); !errors.Is(err, mongo.ErrNoDocuments) {
		utils.WriteError(w, http.StatusBadRequest, "createChat", "chat with these participants exists already")
		return
	}

	if _, err := handler.Models.Chat.Create(participants); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "createChat", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "chat created successfully")
}

func (handler *Handler) GetChat(w http.ResponseWriter, r *http.Request) {
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

	chatId := chi.URLParam(r, "chat_id")
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "chat id is missing")
		return
	}

	chatObjectId, err := utils.ToObjectId(chatId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "strToObjectId", "failed to convert chatId to objectId")
		return
	}

	filter := bson.M{
		"_id": chatObjectId,
		"participants": bson.M{
			"$in": []primitive.ObjectID{userObjectId},
		},
	}

	chatInstance, err2 := handler.Models.Chat.Get(filter, bson.M{})
	if err2 != nil && !errors.Is(err2, mongo.ErrNoDocuments) {
		utils.WriteError(w, http.StatusBadRequest, "getChat", "failed to get chat")
		return
	}

	utils.WriteJSON(w, http.StatusOK, chatInstance)
}

func (handler *Handler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	userObjectId, err := utils.ToObjectId(payload.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "strToObjectId", "failed to convert userId to objectId")
		return
	}

	chatId := chi.URLParam(r, "chat_id")
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "chat id is missing")
		return
	}

	chatObjectId, err := utils.ToObjectId(chatId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "strToObjectId", "failed to convert chatId to objectId")
		return
	}

	filter := bson.M{
		"_id": chatObjectId,
		"participants": bson.M{
			"$in": []primitive.ObjectID{userObjectId},
		},
	}

	if _, err := handler.Models.Chat.Delete(filter); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteChat", "failed to delete chat instance")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "chat deleted successfully")
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
