package handlers

import (
	"chat_app/utils"
	"encoding/hex"
	"errors"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"net/http"
)

func (handler *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "checkAuth", err)
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

	participants := []primitive.ObjectID{payload.UserId, targetUserObjectId}

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

// GetChatMessages -> Returns all the messages of the chat
func (handler *Handler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Type, err.Detail)
		return
	}

	filter := bson.M{
		"sender_id": payload.UserId,
	}

	messages, err2 := handler.Models.Message.GetAll(filter, bson.M{}, 1, 10)
	if err2 != nil {
		utils.WriteError(w, http.StatusBadRequest, "fetchMessages", err2)
		return
	}

	decryptedMessages := make([]string, 0, len(messages))
	for _, message := range messages {
		decodedMessage, err := hex.DecodeString(message.Content)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "msgDecoding", "failed to decode the message")
			continue
		}

		decryptedMsg, err := handler.Cipher.Decrypt(decodedMessage)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "msgDecryption", "failed to decrypt the message")
			continue
		}

		decryptedMessages = append(decryptedMessages, string(decryptedMsg))
	}

	utils.WriteJSON(w, http.StatusOK, decryptedMessages)
}

func (handler *Handler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
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
			"$in": []primitive.ObjectID{payload.UserId},
		},
	}

	if _, err := handler.Models.Chat.Delete(filter); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteChat", "failed to delete chat instance")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "chat deleted successfully")
}

func (handler *Handler) AddChatWebsocket(w http.ResponseWriter, r *http.Request) {
	chatId := chi.URLParam(r, "chat_id")
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "chat id is missing")
		return
	}

	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "user id is missing")
		return
	}

	wsConn, err := WebsocketUpgrade(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "upgradingWebsocket", err)
		return
	}

	wsConn.Add(chatId, userId, handler.WebSocket)

	go func() {
		if err := wsConn.HandleIncomingMessages(chatId, userId, handler.WebSocket, handler); err != nil {
			slog.Error("handling incoming ws messages", "error", err)
		}
	}()
}
