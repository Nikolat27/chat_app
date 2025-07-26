package handlers

import (
	"chat_app/utils"
	"errors"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func (handler *Handler) CreateSaveMessage(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	var input struct {
		Content string `json:"content"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err)
		return
	}

	if _, err := handler.Models.SaveMessage.Create(payload.UserId, "text", input.Content, ""); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "saveMessage", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "msg added to SaveMessages successfully")
}

func (handler *Handler) GetSaveMessages(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"owner_id": payload.UserId,
	}

	msgs, err := handler.Models.SaveMessage.GetAll(filter, bson.M{}, 1, 10)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getMsg", "msg with owner id does not exist")
			return
		}

		utils.WriteError(w, http.StatusBadRequest, "getMsg", "failed to get msg")
		return
	}

	response := map[string]any{
		"messages": msgs,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (handler *Handler) EditSaveMessage(w http.ResponseWriter, r *http.Request) {
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
		"_id":      messageObjectId,
		"owner_id": payload.UserId,
	}

	projection := bson.M{"_id": 1}

	if _, err := handler.Models.SaveMessage.Get(filter, projection); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getMsg", "msg with this id and owner id does not exist")
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

	if _, err := handler.Models.SaveMessage.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "updateMsg", "failed to update the message")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "message updated successfully")
}

func (handler *Handler) DeleteSaveMessage(w http.ResponseWriter, r *http.Request) {
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
		"_id":      messageObjectId,
		"owner_id": payload.UserId,
	}

	if _, err := handler.Models.SaveMessage.Delete(filter); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteMsg", "failed to delete message")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "message deleted successfully")
}
