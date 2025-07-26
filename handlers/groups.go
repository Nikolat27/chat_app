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

func (handler *Handler) RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "group_id")
	if groupId == "" {
		utils.WriteError(w, http.StatusBadRequest, "getParam", "group id is missign")
		return
	}

	groupObjectId, errResp := utils.ToObjectId(groupId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		utils.WriteError(w, http.StatusBadRequest, "getParam", "user id is missign")
		return
	}

	userObjectId, errResp := utils.ToObjectId(userId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"_id": groupObjectId,
	}

	projection := bson.M{
		"_id":      1,
		"users":    1,
		"owner_id": 1,
	}

	if userObjectId == payload.UserId {
		utils.WriteError(w, http.StatusBadRequest, "removeUser", "group owner cant remove himself")
		return
	}

	groupInstance, err := handler.Models.Group.Get(filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getGroup", "group with this id does not exist")
			return
		}

		utils.WriteError(w, http.StatusBadRequest, "getGroup", "failed to get group")
		return
	}

	if groupInstance.OwnerId != payload.UserId {
		utils.WriteError(w, http.StatusBadRequest, "userChecking", "only group owner can remove users")
		return
	}

	if groupInstance.OwnerId == userObjectId {
		utils.WriteError(w, http.StatusBadRequest, "userChecking", "group owner can not remove himself")
		return
	}

	if !slices.Contains(groupInstance.Users, userObjectId) {
		utils.WriteError(w, http.StatusBadRequest, "userChecking", "no user with this id is a member of this group")
		return
	}

	var newUsers []primitive.ObjectID
	for _, user := range groupInstance.Users {
		if user != userObjectId {
			newUsers = append(newUsers, user)
		}
	}

	updates := bson.M{
		"users": newUsers,
	}

	if _, err := handler.Models.Group.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "groupUpdating", "failed to update group")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "user removed successfully")
}
