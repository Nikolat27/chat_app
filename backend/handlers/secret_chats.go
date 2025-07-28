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
		utils.WriteError(w, http.StatusBadRequest, "getSecretChat", "You have a secret chat with this user already")
		return
	}

	if _, err := handler.Models.SecretChat.Create(payload.UserId, targetUserObjectId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "createSecretChat", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "secret chat created successfully")
}

// GetSecretChatMessages -> Returns the whole messages
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
		"chat_id":   chatObjectId,
		"is_secret": true, // secret chat messages
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

		// replace it with the decrypted version
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
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	secretChatId := chi.URLParam(r, "secret_chat_id")
	if secretChatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "secret_chat_id is missing")
		return
	}

	secretChatObjectId, errResp := utils.ToObjectId(secretChatId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, "strToObjectId", "invalid secret_chat_id")
		return
	}

	var input struct {
		PublicKey string `json:"public_key"`
	}
	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err)
		return
	}

	if input.PublicKey == "" {
		utils.WriteError(w, http.StatusBadRequest, "missingPublicKey", "public key is required")
		return
	}

	// Fetch the secret chat
	filter := bson.M{"_id": secretChatObjectId}
	secretChat, err := handler.Models.SecretChat.Get(filter, bson.M{})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "getSecretChat", err)
		return
	}

	// Check which user is making the request
	isUser1 := secretChat.User1 == payload.UserId
	isUser2 := secretChat.User2 == payload.UserId

	if !isUser1 && !isUser2 {
		utils.WriteError(w, http.StatusForbidden, "unAuthorized", "You are not a participant in this secret chat")
		return
	}

	// Prevent overriding existing keys
	if isUser1 && secretChat.User1PublicKey != "" {
		utils.WriteError(w, http.StatusConflict, "keyExists", "User 1 public key already set")
		return
	}
	if isUser2 && secretChat.User2PublicKey != "" {
		utils.WriteError(w, http.StatusConflict, "keyExists", "User 2 public key already set")
		return
	}

	// Prepare update
	updates := bson.M{}
	if isUser1 {
		updates["user_1_public_key"] = input.PublicKey
	}
	if isUser2 {
		updates["user_2_public_key"] = input.PublicKey
	}

	if _, err := handler.Models.SecretChat.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "updateSecretChat", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Public key saved successfully")
}

func (handler *Handler) ApproveSecretChat(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	chatID := chi.URLParam(r, "secret_chat_id")

	objectId, _ := primitive.ObjectIDFromHex(chatID)

	filter := bson.M{
		"_id":    objectId,
		"user_2": payload.UserId,
	}

	update := bson.M{
		"user_2_accepted": true,
	}

	if _, err := handler.Models.SecretChat.Update(filter, update); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "approveSecretChat", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "secret chat approved")
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
