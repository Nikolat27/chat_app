package handlers

import (
	"chat_app/utils"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
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

	file, header, err2 := r.FormFile("file")
	if err2 != nil {
		utils.WriteError(w, http.StatusBadRequest, "getFile", "failed to get form file")
		return
	}
	defer file.Close()

	if err := os.MkdirAll("uploads", 0755); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "MkdirAll", "failed to create upload dir")
		return
	}

	fileName := rand.Text() + header.Filename
	path := filepath.Join("uploads", fileName)

	dst, err2 := os.Create(path)
	if err2 != nil {
		utils.WriteError(w, http.StatusBadRequest, "createPath", err2)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "saveFile", "failed to save file")
		return
	}

	senderId := payload.UserId
	if _, err := handler.Models.Message.Create(chatObjectId, senderId, receiverObjectId, "image",
		fileName, ""); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "createMsg", "failed to create message")
		return
	}

	resp := map[string]string{
		"url": fileName,
	}

	utils.WriteJSON(w, http.StatusCreated, resp)
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

	wsConn.Add(chatId, senderId, handler.WebSocket)

	go func() {
		if err := wsConn.HandleIncomingMessages(chatId, senderId, receiverId, handler.WebSocket, handler); err != nil {
			slog.Error("handling incoming ws messages", "error", err)
		}
	}()
}
