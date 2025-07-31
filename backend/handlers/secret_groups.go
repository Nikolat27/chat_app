package handlers

import (
	"chat_app/utils"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

func (handler *Handler) CreateSecretGroup(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
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
	groupType := r.FormValue("group_type")
	if groupType == "" {
		groupType = "public"
	}
	publicKey := r.FormValue("public_key")
	if publicKey == "" {
		utils.WriteError(w, http.StatusBadRequest, "missingField", "public key is required")
		return
	}

	inviteLink := uuid.New().String()
	members := []primitive.ObjectID{payload.UserId}
	admins := []primitive.ObjectID{payload.UserId}
	joinTimes := map[string]time.Time{payload.UserId.Hex(): time.Now()}

	result, err := handler.Models.SecretGroup.Create(
		payload.UserId, name, description, groupType, inviteLink,
		members, admins, joinTimes,
	)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "createGroup", "failed to create group")
		return
	}

	resp := map[string]string{
		"message":     "secret group created successfully",
		"group_id":    result.InsertedID.(primitive.ObjectID).Hex(),
		"owner_id":    payload.UserId.Hex(),
		"invite_link": inviteLink,
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func (handler *Handler) UpdateSecretGroup(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "secret_group_id")
	if groupId == "" {
		utils.WriteError(w, http.StatusUnauthorized, "getUrlParam", "group id is missing")
		return
	}

	groupObjectId, errResp := utils.ToObjectId(groupId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		utils.WriteError(w, http.StatusUnauthorized, "formValue", "name field is empty")
		return
	}

	description := r.FormValue("description")
	if description == "" {
		utils.WriteError(w, http.StatusUnauthorized, "formValue", "description field is empty")
		return
	}

	groupType := r.FormValue("group_type")
	if groupType == "" {
		utils.WriteError(w, http.StatusUnauthorized, "formValue", "group_type field is empty")
		return
	}

	var inviteLink string
	if groupType == "private" {
		inviteLink = uuid.New().String()
	}

	filter := bson.M{
		"_id":      groupObjectId,
		"owner_id": payload.UserId,
	}

	updates := bson.M{
		"name":        name,
		"description": description,
		"group_type":  groupType,
		"invite_link": inviteLink,
	}

	_, err := handler.Models.SecretGroup.Update(filter, updates)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "updateGroup", "failed to update group")
		return
	}

	response := map[string]string{
		"invite_link": inviteLink,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (handler *Handler) JoinSecretGroup(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
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
		"_id": 1, "members": 1, "banned_members": 1, "type": 1,
	}

	groupInstance, err := handler.Models.SecretGroup.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "getGroup", "group with this invite link does not exist")
		return
	}

	if slices.Contains(groupInstance.BannedMembers, payload.UserId) {
		utils.WriteError(w, http.StatusBadRequest, "userBanned", "you are banned from this group")
		return
	}

	if slices.Contains(groupInstance.Members, payload.UserId) {
		utils.WriteError(w, http.StatusBadRequest, "userExists", "you are already in this group")
		return
	}

	if groupInstance.Type == "private" {
		if err := checkUserApproval(groupInstance.Id, payload.UserId, handler); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
			return
		}
	}

	filter = bson.M{
		"_id": groupInstance.Id,
	}

	update := bson.M{
		"$push": bson.M{
			"members": payload.UserId,
		},
		"$set": bson.M{
			"member_join_times." + payload.UserId.Hex(): time.Now(),
		},
	}

	if _, err := handler.Models.SecretGroup.UpdateSmart(filter, update); err != nil {
		fmt.Println(err)
		utils.WriteError(w, http.StatusBadRequest, "joinGroup", "failed to join group")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "user joined secret group successfully")
}

func (handler *Handler) RemoveUserFromSecretGroup(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "secret_group_id")
	if groupId == "" {
		utils.WriteError(w, http.StatusBadRequest, "getParam", "group id is missing")
		return
	}

	groupObjectId, errResp := utils.ToObjectId(groupId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		utils.WriteError(w, http.StatusBadRequest, "getParam", "user id is missing")
		return
	}

	targetUserObjectId, errResp := utils.ToObjectId(userId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{"_id": groupObjectId}
	projection := bson.M{"owner_id": 1, "members": 1}

	groupInstance, err := handler.Models.SecretGroup.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "getGroup", err.Error())
		return
	}

	if payload.UserId != groupInstance.OwnerId {
		utils.WriteError(w, http.StatusForbidden, "auth", "only owner can remove users")
		return
	}

	if groupInstance.OwnerId == targetUserObjectId {
		utils.WriteError(w, http.StatusBadRequest, "auth", "owner can't be removed")
		return
	}

	newMembers := utils.DeleteElementFromSlice(groupInstance.Members, targetUserObjectId)
	updates := bson.M{"members": newMembers}
	if _, err := handler.Models.SecretGroup.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "update", err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, "user removed successfully")
}

func (handler *Handler) DeleteSecretGroup(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "secret_group_id")
	groupObjectId, errResp := utils.ToObjectId(groupId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{"_id": groupObjectId, "owner_id": payload.UserId}
	if _, err := handler.Models.SecretGroup.Delete(filter); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "delete", err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, "group deleted successfully")
}

func (handler *Handler) LeaveSecretGroup(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "secret_group_id")
	groupObjectId, errResp := utils.ToObjectId(groupId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{"_id": groupObjectId}
	projection := bson.M{"owner_id": 1, "members": 1}
	groupInstance, err := handler.Models.SecretGroup.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "getGroup", err.Error())
		return
	}

	if payload.UserId == groupInstance.OwnerId {
		utils.WriteError(w, http.StatusBadRequest, "leave", "owner can't leave the group")
		return
	}

	newMembers := utils.DeleteElementFromSlice(groupInstance.Members, payload.UserId)
	updates := bson.M{"members": newMembers}
	if _, err := handler.Models.SecretGroup.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "update", err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, "left group successfully")
}

func (handler *Handler) getSecretGroupIdByInviteLink(inviteLink string) (primitive.ObjectID, error) {
	filter := bson.M{"invite_link": inviteLink}
	projection := bson.M{"_id": 1}
	groupInstance, err := handler.Models.SecretGroup.Get(filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return primitive.NilObjectID, errors.New("group with this id does not exist")
		}
		return primitive.NilObjectID, errors.New(err.Error())
	}
	return groupInstance.Id, nil
}

// BanMemberFromSecretGroup -> Ban users
func (handler *Handler) BanMemberFromSecretGroup(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "secret_group_id")
	if groupId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "group id is missing")
		return
	}

	groupObjectId, errResp := utils.ToObjectId(groupId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	var input struct {
		TargetUser string `json:"target_user"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err.Error())
		return
	}

	targetUserObjectId, errResp := utils.ToObjectId(input.TargetUser)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{"_id": groupObjectId}
	group, err := handler.Models.SecretGroup.Get(filter, bson.M{})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "getGroup", err.Error())
		return
	}

	if !slices.Contains(group.Admins, payload.UserId) {
		utils.WriteError(w, http.StatusForbidden, "unauthorized", "Only admins can ban users")
		return
	}

	if slices.Contains(group.BannedMembers, targetUserObjectId) {
		utils.WriteError(w, http.StatusBadRequest, "alreadyBanned", "User already banned")
		return
	}

	group.Members = utils.DeleteElementFromSlice(group.Members, targetUserObjectId)
	group.BannedMembers = append(group.BannedMembers, targetUserObjectId)

	updates := bson.M{
		"members":        group.Members,
		"banned_members": group.BannedMembers,
	}

	if _, err := handler.Models.SecretGroup.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "updateGroup", err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, "user banned successfully")
}

// UnBanMemberFromSecretGroup -> UnBan users
func (handler *Handler) UnBanMemberFromSecretGroup(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "secret_group_id")
	if groupId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "group id is missing")
		return
	}

	groupObjectId, errResp := utils.ToObjectId(groupId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	var input struct {
		TargetUser string `json:"target_user"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err.Error())
		return
	}

	targetUserObjectId, errResp := utils.ToObjectId(input.TargetUser)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{"_id": groupObjectId}
	group, err := handler.Models.SecretGroup.Get(filter, bson.M{})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "getGroup", err.Error())
		return
	}

	if !slices.Contains(group.Admins, payload.UserId) {
		utils.WriteError(w, http.StatusForbidden, "unauthorized", "Only admins can unban users")
		return
	}

	if !slices.Contains(group.BannedMembers, targetUserObjectId) {
		utils.WriteError(w, http.StatusBadRequest, "notBanned", "User is not banned")
		return
	}

	group.BannedMembers = utils.DeleteElementFromSlice(group.BannedMembers, targetUserObjectId)
	group.Members = append(group.Members, targetUserObjectId)

	updates := bson.M{
		"banned_members": group.BannedMembers,
		"members":        group.Members,
	}

	if _, err := handler.Models.SecretGroup.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "updateGroup", err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, "user unbanned successfully")
}

func (handler *Handler) GetSecretGroupMembers(w http.ResponseWriter, r *http.Request) {
	_, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "secret_group_id")
	if groupId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "group id is missing")
		return
	}

	groupObjectId, errResp := utils.ToObjectId(groupId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{"_id": groupObjectId}
	projection := bson.M{"members": 1, "user_public_keys": 1}
	groupInstance, err := handler.Models.SecretGroup.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "fetchGroup", err)
		return
	}

	members := make(map[string]map[string]string)
	for _, userId := range groupInstance.Members {
		username, _ := getUserUsername(userId, handler)
		avatarUrl, _ := getUserAvatarUrl(userId, handler)

		members[userId.Hex()] = map[string]string{
			"username":   username,
			"avatar_url": avatarUrl,
		}
	}

	utils.WriteJSON(w, http.StatusOK, members)
}

func (handler *Handler) GetSecretGroupMessages(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "secret_group_id")
	if groupId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "group id is missing")
		return
	}

	groupObjectId, errResp := utils.ToObjectId(groupId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	groupFilter := bson.M{"_id": groupObjectId}
	projection := bson.M{"members": 1}

	groupInstance, err := handler.Models.SecretGroup.Get(groupFilter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "getGroup", "failed to fetch group")
		return
	}

	if !slices.Contains(groupInstance.Members, payload.UserId) {
		utils.WriteError(w, http.StatusBadRequest, "getMsgs", "you are not a member of this group")
		return
	}

	page, pageLimit, errResp := utils.ParsePageAndLimitQueryParams(r.URL)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"group_id": groupObjectId,
	}

	messages, err := handler.Models.SecretGroupMessages.GetAll(filter, bson.M{}, page, pageLimit)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "fetchMessages", err)
		return
	}

	resp := map[string]any{
		"messages": messages,
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

// AddSecretGroupWebsocket -> Establish Websocket and send the messages
func (handler *Handler) AddSecretGroupWebsocket(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "secret_group_id")
	if groupId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "group id is missing")
		return
	}

	senderId := r.URL.Query().Get("sender_id")
	if senderId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "sender id is missing")
		return
	}

	wsConn, err := WebsocketUpgrade(w, r)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "upgradingWebsocket", err)
		return
	}

	wsConn.AddGroup(groupId, senderId, handler.WebSocket)

	go func() {
		if err := wsConn.HandleSecretGroupIncomingMsgs(groupId, senderId, handler.WebSocket, handler); err != nil {
			slog.Error("handling incoming ws messages", "error", err)
		}
	}()
}
