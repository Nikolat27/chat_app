package handlers

import (
	"chat_app/utils"
	"encoding/hex"
	"errors"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func (handler *Handler) storeChatMsgToDB(chatId, senderId, receiverId string, isSecret bool, payload []byte) error {
	chatObjectId, err := utils.ToObjectId(chatId)
	if err != nil {
		return errors.New(err.Type)
	}

	senderObjectId, err := utils.ToObjectId(senderId)
	if err != nil {
		return errors.New(err.Type)
	}

	receiverObjectId, err := utils.ToObjectId(receiverId)
	if err != nil {
		return errors.New(err.Type)
	}

	ciphered, err2 := handler.Cipher.Encrypt(payload)
	if err2 != nil {
		return err2
	}

	encodedCipher := hex.EncodeToString(ciphered)

	if _, err := handler.Models.Message.Create(chatObjectId, primitive.NilObjectID, senderObjectId, receiverObjectId,
		"text", "", encodedCipher, isSecret); err != nil {
		return err
	}

	return nil
}

func (handler *Handler) storeGroupMsgToDB(groupId, senderId string, isSecret bool, payload []byte) error {
	senderObjectId, err := utils.ToObjectId(senderId)
	if err != nil {
		return errors.New(err.Type)
	}

	ciphered, err2 := handler.Cipher.Encrypt(payload)
	if err2 != nil {
		return err2
	}

	encodedCipher := hex.EncodeToString(ciphered)

	groupObjectId, err := utils.ToObjectId(groupId)
	if err != nil {
		return errors.New(err.Type)
	}

	if _, err := handler.Models.Message.Create(primitive.NilObjectID, groupObjectId, senderObjectId, primitive.NilObjectID,
		"text", "", encodedCipher, isSecret); err != nil {
		return err
	}

	return nil
}

func (handler *Handler) storeSecretGroupMsgToDB(groupId, senderId, payload string) error {
	senderObjectId, err := utils.ToObjectId(senderId)
	if err != nil {
		return errors.New(err.Type)
	}

	groupObjectId, err := utils.ToObjectId(groupId)
	if err != nil {
		return errors.New(err.Type)
	}

	if _, err := handler.Models.SecretGroupMessages.Create(groupObjectId, senderObjectId, payload); err != nil {
		return err
	}

	return nil
}

func (handler *Handler) EditMessage(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	messageId := chi.URLParam(r, "message_id")
	if messageId == "" {
		utils.WriteError(w, http.StatusBadRequest, "missingParam", "message id is missing")
		return
	}

	messageObjectId, err := utils.ToObjectId(messageId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	filter := bson.M{
		"_id":       messageObjectId,
		"sender_id": payload.UserId,
	}

	projection := bson.M{"_id": 1}

	if _, err := handler.Models.Message.Get(filter, projection); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getMsg", "msg with this id and sender id does not exist")
			return
		}

		utils.WriteError(w, http.StatusBadRequest, "getMsg", "failed to get msg")
		return
	}

	var input struct {
		NewContent string `json:"new_content"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", "failed to parse the body data")
		return
	}

	updates := bson.M{
		"content": input.NewContent,
	}

	if _, err := handler.Models.Message.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "updateMsg", "failed to update the message")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "message updated successfully")
}

func (handler *Handler) DeleteMessageForSender(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	messageId := chi.URLParam(r, "message_id")
	if messageId == "" {
		utils.WriteError(w, http.StatusBadRequest, "missingParam", "message id is missing")
		return
	}

	messageObjectId, err := utils.ToObjectId(messageId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	filter := bson.M{
		"_id":       messageObjectId,
		"sender_id": payload.UserId,
	}

	projection := bson.M{"_id": 1}

	if _, err := handler.Models.Message.Get(filter, projection); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getMsg", "msg with this id and sender id does not exist")
			return
		}

		utils.WriteError(w, http.StatusBadRequest, "getMsg", "failed to get msg")
		return
	}

	filter = bson.M{
		"_id": messageObjectId,
	}

	updates := bson.M{
		"is_deleted_for_sender": true,
	}

	if _, err := handler.Models.Message.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteMsg", "failed to delete msg")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "message deleted successfully")
}

func (handler *Handler) DeleteMessageForReceiver(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	messageId := chi.URLParam(r, "message_id")
	if messageId == "" {
		utils.WriteError(w, http.StatusBadRequest, "missingParam", "message id is missing")
		return
	}

	messageObjectId, err := utils.ToObjectId(messageId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	filter := bson.M{
		"_id":         messageObjectId,
		"receiver_id": payload.UserId,
	}

	projection := bson.M{
		"_id": 1,
	}

	if _, err := handler.Models.Message.Get(filter, projection); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getMsg", "msg with this id and receiver id does not exist")
			return
		}

		utils.WriteError(w, http.StatusBadRequest, "getMsg", "failed to get msg")
		return
	}

	filter = bson.M{
		"_id": messageObjectId,
	}

	updates := bson.M{
		"is_deleted_for_receiver": true,
	}

	if _, err := handler.Models.Message.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteMsg", "failed to delete msg")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "message deleted successfully")
}

func (handler *Handler) DeleteMessageForAll(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	messageId := chi.URLParam(r, "message_id")
	if messageId == "" {
		utils.WriteError(w, http.StatusBadRequest, "missingParam", "message id is missing")
		return
	}

	messageObjectId, err := utils.ToObjectId(messageId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	filter := bson.M{
		"_id": messageObjectId,
		"$or": []bson.M{
			{"sender_id": payload.UserId},
			{"receiver_id": payload.UserId},
		},
	}

	deletedResult, err2 := handler.Models.Message.Delete(filter)
	if err2 != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteMsg", "failed to delete msg")
		return
	}

	if deletedResult.DeletedCount == 0 {
		utils.WriteError(w, http.StatusBadRequest, "deleteMsg", "no msg with this id exists")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "message deleted successfully")
}
