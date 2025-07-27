package handlers

import (
	"chat_app/utils"
	"encoding/hex"
	"errors"
	"fmt"
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

func (handler *Handler) UploadChatImage(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "checkAuth", err)
		return
	}

	chatId := chi.URLParam(r, "chat_id")
	if chatId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "chat id is missing")
		return
	}

	chatObjectId, err := utils.ToObjectId(chatId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "strToObjectId", "failed to convert chat id to objectId")
		return
	}

	receiverId := chi.URLParam(r, "receiver_id")
	if receiverId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "receiver id is missing")
		return
	}

	receiverObjectId, err := utils.ToObjectId(receiverId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "strToObjectId", "failed to convert chat id to objectId")
		return
	}

	filter := bson.M{
		"_id": chatObjectId,
		"participants": bson.M{
			"$in": []primitive.ObjectID{payload.UserId},
		},
	}

	projection := bson.M{
		"_id": 1,
	}

	if _, err := handler.Models.Chat.Get(filter, projection); err != nil {
		fmt.Println(err)
		utils.WriteError(w, http.StatusBadRequest, "getChat", "failed to get chat instance")
		return
	}

	allowedFormats := []string{".jpg", ".jpeg", ".png", ".webp"}
	avatarAddress, err := utils.UploadFile(r, "file", allowedFormats)

	senderId := payload.UserId
	if _, err := handler.Models.Message.Create(chatObjectId, primitive.NilObjectID, senderId, receiverObjectId, "image",
		avatarAddress, ""); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "createMsg", "failed to create message")
		return
	}

	resp := map[string]string{
		"url": avatarAddress,
	}

	utils.WriteJSON(w, http.StatusCreated, resp)
}

// GetChatMessages -> Returns all the messages of the chat
func (handler *Handler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	if _, errResp := utils.CheckAuth(r.Header, handler.Paseto); errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	chatId := chi.URLParam(r, "chat_id")
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
		"chat_id": chatObjectId,
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
		if err := wsConn.HandleChatIncomingMsgs(chatId, senderId, receiverId, handler.WebSocket, handler); err != nil {
			slog.Error("handling incoming ws messages", "error", err)
		}
	}()
}
