package handlers

import (
	"chat_app/utils"
	"crypto/rand"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"slices"
)

func (handler *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Type, err.Detail)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		name = rand.Text()
	}

	description := r.FormValue("description")
	if description == "" {
		description = "no description"
	}

	allowedFormats := []string{".jpg", ".jpeg", ".png", ".webp"}
	avatarUrl, err := utils.UploadFile(r, "file", allowedFormats)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	inviteLink := uuid.New().String()

	users := []primitive.ObjectID{payload.UserId}

	if _, err := handler.Models.Group.Create(payload.UserId, name, description, avatarUrl, inviteLink, users); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "createGroup", "failed to create group")
		return
	}

	response := map[string]string{
		"message":     "group created successfully",
		"invite_link": inviteLink,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (handler *Handler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Type, err.Detail)
		return
	}

	inviteLink := chi.URLParam(r, "invite_link")
	if inviteLink == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "invite link is missing")
		return
	}

	filter := bson.M{
		"invite_link": inviteLink,
	}

	projection := bson.M{
		"_id":   1,
		"users": 1,
	}

	groupInstance, err2 := handler.Models.Group.Get(filter, projection)
	if err2 != nil {
		if errors.Is(err2, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getGroup", "group with this invite link does not exist")
			return
		}

		utils.WriteError(w, http.StatusBadRequest, "getGroup", "failed to get group")
		return
	}

	if slices.Contains(groupInstance.Users, payload.UserId) {
		utils.WriteError(w, http.StatusBadRequest, "userExists", "This user is already in the group")
		return
	}

	newUsers := append(groupInstance.Users, payload.UserId)

	updates := bson.M{
		"users": newUsers,
	}

	if _, err := handler.Models.Group.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "updateGroup", "failed to join the users")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "user joined successfully")
}
