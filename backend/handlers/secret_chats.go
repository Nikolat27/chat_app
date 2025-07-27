package handlers

import (
	"chat_app/database/models"
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

func (handler *Handler) CreateSecretChat(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	var input struct {
		TargetUser string `json:"target_user"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err)
		return
	}

	targetUserObjectId, errResp := utils.ToObjectId(input.TargetUser)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"$or": []bson.M{
			{"user_1": payload.UserId, "user_2": targetUserObjectId},
			{"user_1": targetUserObjectId, "user_2": payload.UserId},
		},
	}

	projection := bson.M{
		"_id": 1,
	}

	_, err := handler.Models.SecretChat.Get(filter, projection)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		utils.WriteError(w, http.StatusBadRequest, "getSecretChat", err)
		return
	} else if err == nil {
		utils.WriteError(w, http.StatusBadRequest, "getSecretChat", "You have a secret chat with these users already")
		return
	}

	if _, err := handler.Models.SecretChat.Create(payload.UserId, targetUserObjectId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "createSecretChat", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "secret chat created successfully")
}

func (handler *Handler) GetSecretChatMessages(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	chatId := chi.URLParam(r, "secret_chat_id")
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "chat id is missing")
		return
	}

	chatObjectId, errResp := utils.ToObjectId(chatId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, "strToObjectId", "failed to convert chatId")
		return
	}

	filter := bson.M{
		"chat_id": chatObjectId, "is_secret": true,
	}

	page, limit, errResp := utils.ParsePageAndLimitQueryParams(r.URL)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	messages, err := handler.Models.Message.GetAll(filter, bson.M{}, page, limit)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "fetchMessages", err)
		return
	}

	for idx := range messages {
		if messages[idx].IsDeletedForSender && messages[idx].SenderId == payload.UserId {
			messages[idx] = models.Message{} // skip it
			continue
		}
		decodedMessage, err := hex.DecodeString(messages[idx].Content)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "msgDecoding", "failed to decode the message")
			continue
		}
		decryptedMsg, err := handler.Cipher.Decrypt(decodedMessage)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "msgDecryption", "failed to decrypt the message")
			continue
		}
		messages[idx].Content = string(decryptedMsg)
	}

	utils.WriteJSON(w, http.StatusOK, messages)
}

func (handler *Handler) DeleteSecretChat(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	chatId := chi.URLParam(r, "secret_chat_id")
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "chat id is missing")
		return
	}

	chatObjectId, err := utils.ToObjectId(chatId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "strToObjectId", "failed to convert chatId")
		return
	}

	filter := bson.M{
		"_id": chatObjectId,
		"participants": bson.M{
			"$in": []primitive.ObjectID{payload.UserId},
		},
	}

	if _, err := handler.Models.SecretChat.Delete(filter); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteSecretChat", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "secret chat deleted successfully")
}

func (handler *Handler) UpdateSecretChat(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.CheckAuth(r.Header, handler.Paseto); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	secretChatId := chi.URLParam(r, "secret_chat_id")
	if secretChatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "chat id is missing")
		return
	}

	secretChatObjectId, err := utils.ToObjectId(secretChatId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "strToObjectId", "failed to convert secretChatId")
		return
	}

	var input struct {
		User2Accepted  bool   `json:"user_2_accepted"`
		User1PublicKey string `json:"user_1_public_key"`
		User2PublicKey string `json:"user_2_public_key"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err)
		return
	}

	filter := bson.M{
		"_id": secretChatObjectId,
	}

	updates := bson.M{
		"user_2_accepted":   input.User2Accepted,
		"user_1_public_key": input.User1PublicKey,
		"user_2_public_key": input.User2PublicKey,
	}

	if _, err := handler.Models.SecretChat.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteSecretChat", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "secret chat updated successfully")
}

func (handler *Handler) AddSecretChatWebsocket(w http.ResponseWriter, r *http.Request) {
	chatId := chi.URLParam(r, "secret_chat_id")
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "chat id is missing")
		return
	}

	senderId := r.URL.Query().Get("sender_id")
	if senderId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "sender id is missing")
		return
	}

	receiverId := r.URL.Query().Get("receiver_id")
	if receiverId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "receiver id is missing")
		return
	}

	wsConn, err := WebsocketUpgrade(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "upgradingWebsocket", err)
		return
	}

	wsConn.AddChat(chatId, senderId, handler.WebSocket)

	go func() {
		if err := wsConn.HandleChatIncomingMsgs(chatId, senderId, receiverId, true, handler.WebSocket, handler); err != nil {
			slog.Error("handling incoming secret chat ws messages", "error", err)
		}
	}()
}
