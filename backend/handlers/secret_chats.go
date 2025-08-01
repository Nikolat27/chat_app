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

func (handler *Handler) GetSecretChat(w http.ResponseWriter, r *http.Request) {
	_, errResp := utils.CheckAuth(r, handler.Paseto)
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
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"_id": chatObjectId,
	}

	chatInstance, err := handler.Models.SecretChat.Get(filter, bson.M{})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "getSecretChat", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, chatInstance)
}

func (handler *Handler) CreateSecretChat(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r, handler.Paseto)
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

	result, err := handler.Models.SecretChat.Create(payload.UserId, targetUserObjectId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "createSecretChat", err)
		return
	}

	resp := map[string]string{
		"secret_chat_id": result.Hex(),
	}

	utils.WriteJSON(w, http.StatusCreated, resp)
}

// GetSecretChatMessages -> Returns the whole messages
func (handler *Handler) GetSecretChatMessages(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r, handler.Paseto)
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
	payload, errResp := utils.CheckAuth(r, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	chatId := chi.URLParam(r, "secret_chat_id")
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "chat id is missing")
		return
	}

	chatObjectId, errResp := utils.ToObjectId(chatId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"_id": chatObjectId,
		"$or": []bson.M{
			{"user_1": payload.UserId},
			{"user_2": payload.UserId},
		},
	}

	if _, err := handler.Models.SecretChat.Delete(filter); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteSecretChat", err)
		return
	}

	filter = bson.M{
		"chat_id":   chatObjectId,
		"is_secret": true,
	}

	if _, err := handler.Models.Message.DeleteAll(filter); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteSecretChatMessages", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "secret chat deleted successfully + its messages")
}

func (handler *Handler) UploadSecretChatPublicKey(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r, handler.Paseto)
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
		utils.WriteError(w, http.StatusBadRequest, "missingKey", "public key is required")
		return
	}

	// Fetch chat
	filter := bson.M{"_id": secretChatObjectId}
	secretChat, err := handler.Models.SecretChat.Get(filter, bson.M{})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "getSecretChat", err)
		return
	}

	isUser1 := secretChat.User1 == payload.UserId
	isUser2 := secretChat.User2 == payload.UserId
	if !isUser1 && !isUser2 {
		utils.WriteError(w, http.StatusForbidden, "unauthorized", "You're not part of this chat")
		return
	}

	if isUser1 && secretChat.User1PublicKey != "" {
		utils.WriteError(w, http.StatusConflict, "keyExists", "User 1 public key already set")
		return
	}
	if isUser2 && secretChat.User2PublicKey != "" {
		utils.WriteError(w, http.StatusConflict, "keyExists", "User 2 public key already set")
		return
	}

	updates := bson.M{}
	if isUser1 {
		updates["user_1_public_key"] = input.PublicKey
	}
	if isUser2 {
		updates["user_2_public_key"] = input.PublicKey
	}

	if _, err := handler.Models.SecretChat.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "updatePublicKey", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Public key saved successfully")
}

func (handler *Handler) UploadSecretChatSymmetricKey(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r, handler.Paseto)
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
		User1EncryptedSymmetricKey string `json:"user_1_encrypted_symmetric_key"`
		User2EncryptedSymmetricKey string `json:"user_2_encrypted_symmetric_key"`
	}

	if err := utils.ParseJSON(r.Body, 1000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err)
		return
	}

	if input.User1EncryptedSymmetricKey == "" || input.User2EncryptedSymmetricKey == "" {
		utils.WriteError(w, http.StatusBadRequest, "missingKey", "encrypted symmetric key is missing")
		return
	}

	// Fetch chat
	filter := bson.M{
		"_id": secretChatObjectId,
		"$or": []bson.M{
			{"user_1": payload.UserId},
			{"user_2": payload.UserId},
		},
	}

	secretChat, err := handler.Models.SecretChat.Get(filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getSecretChat", "secret chat with these users or Id does not exist")
			return
		}

		utils.WriteError(w, http.StatusBadRequest, "getSecretChat", "failed to get the secret chat")
		return
	}

	// Check that both public keys exist before saving the chat encrypted key
	if secretChat.User1PublicKey == "" || secretChat.User2PublicKey == "" {
		utils.WriteError(w, http.StatusBadRequest, "missingPublicKeys", "Both public keys must be uploaded before setting encrypted symmetric key")
		return
	}

	updates := bson.M{
		"user_1_encrypted_symmetric_key": input.User1EncryptedSymmetricKey,
		"user_2_encrypted_symmetric_key": input.User2EncryptedSymmetricKey,
		"key_finalized":                  true,
	}

	if _, err := handler.Models.SecretChat.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "updateSymmetricKey", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Encrypted symmetric key saved successfully")
}

func (handler *Handler) ApproveSecretChat(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r, handler.Paseto)
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
