package handlers

import (
	"chat_app/utils"
	"crypto/rand"
	"encoding/hex"
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
)

func (handler *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
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

	allowedFormats := []string{".jpg", ".jpeg", ".png", ".webp"}
	avatarUrl, errResp := utils.UploadFile(r, "file", allowedFormats)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	inviteLink := uuid.New().String()

	members := []primitive.ObjectID{payload.UserId}

	admins := []primitive.ObjectID{payload.UserId}

	result, err := handler.Models.Group.Create(payload.UserId, name, description, avatarUrl, groupType, inviteLink, members, admins)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "createGroup", "failed to create group")
		return
	}

	response := map[string]string{
		"message":     "group created successfully",
		"group_id":    result.InsertedID.(primitive.ObjectID).Hex(),
		"owner_id":    payload.UserId.Hex(),
		"invite_link": inviteLink,
		"avatar_url":  avatarUrl,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (handler *Handler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "group_id")
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

	allowedFormats := []string{".jpg", ".jpeg", ".png", ".webp"}
	avatarUrl, errResp := utils.UploadFile(r, "file", allowedFormats)
	if errResp != nil {
		if errResp.Type == "fileMissing" {
			avatarUrl = ""
		} else {
			utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
			return
		}
	}

	var inviteLink string
	if groupType == "private" {
		inviteLink = uuid.New().String()
	}

	filter := bson.M{
		"_id":      groupObjectId,
		"owner_id": payload.UserId,
	}

	var updates bson.M
	if avatarUrl != "" {
		updates = bson.M{
			"name":        name,
			"description": description,
			"group_type":  groupType,
			"invite_link": inviteLink,
			"avatar_url":  avatarUrl,
		}
	} else {
		updates = bson.M{
			"name":        name,
			"description": description,
			"group_type":  groupType,
			"invite_link": inviteLink,
		}
	}

	_, err := handler.Models.Group.Update(filter, updates)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "createGroup", "failed to create group")
		return
	}

	response := map[string]string{
		"invite_link": inviteLink,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (handler *Handler) JoinGroup(w http.ResponseWriter, r *http.Request) {
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
		"_id":            1,
		"members":        1,
		"banned_members": 1,
		"type":           1,
	}

	groupInstance, err := handler.Models.Group.Get(filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getGroup", "group with this invite link does not exist")
			return
		}

		utils.WriteError(w, http.StatusBadRequest, "getGroup", "failed to get group")
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

	newMembers := append(groupInstance.Members, payload.UserId)

	updates := bson.M{
		"members": newMembers,
	}

	if _, err := handler.Models.Group.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "updateGroup", "failed to join the members")
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
		utils.WriteError(w, http.StatusBadRequest, "getParam", "member id is missign")
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
		"members":  1,
		"owner_id": 1,
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
		utils.WriteError(w, http.StatusBadRequest, "userChecking", "only group owner can remove members")
		return
	}

	if groupInstance.OwnerId == userObjectId {
		utils.WriteError(w, http.StatusBadRequest, "userChecking", "group owner can not remove himself")
		return
	}

	if !slices.Contains(groupInstance.Members, userObjectId) {
		utils.WriteError(w, http.StatusBadRequest, "userChecking", "no member with this id is a member of this group")
		return
	}

	var newMembers []primitive.ObjectID
	for _, member := range groupInstance.Members {
		if member != userObjectId {
			newMembers = append(newMembers, member)
		}
	}

	updates := bson.M{
		"members": newMembers,
	}

	if _, err := handler.Models.Group.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "groupUpdating", "failed to update group")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "member removed successfully")
}

func (handler *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
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

	filter := bson.M{
		"_id": groupObjectId,
	}

	projection := bson.M{
		"_id":      1,
		"members":  1,
		"owner_id": 1,
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
		utils.WriteError(w, http.StatusBadRequest, "userChecking", "only group owner can delete it")
		return
	}

	if _, err := handler.Models.Group.Delete(filter); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteGroup", "failed to delete group")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "group deleted successfully")
}

func (handler *Handler) AddGroupWebsocket(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "group_id")
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
		if err := wsConn.HandleGroupIncomingMsgs(groupId, senderId, false, handler.WebSocket, handler); err != nil {
			slog.Error("handling incoming ws messages", "error", err)
		}
	}()
}

// GetGroupMessages -> Returns all the messages of the group
func (handler *Handler) GetGroupMessages(w http.ResponseWriter, r *http.Request) {
	if _, errResp := utils.CheckAuth(r.Header, handler.Paseto); errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "group_id")
	if groupId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "group id is missing")
		return
	}

	groupObjectId, errResp := utils.ToObjectId(groupId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"group_id": groupObjectId,
	}

	page, pageLimit, errResp := utils.ParsePageAndLimitQueryParams(r.URL)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	messages, err := handler.Models.Message.GetAll(filter, bson.M{}, page, pageLimit)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "fetchMessages", err)
		return
	}

	for idx := range messages {
		decodedMessage, err := hex.DecodeString(messages[idx].Content)
		if err != nil {
			slog.Warn("failed to decode message", "err", err, "msgID", messages[idx].Id.Hex())
			continue
		}

		decryptedMsg, err := handler.Cipher.Decrypt(decodedMessage)
		if err != nil {
			slog.Warn("failed to decrypt message", "err", err, "msgID", messages[idx].Id.Hex())
			continue
		}

		messages[idx].Content = string(decryptedMsg)
	}

	resp := map[string]any{
		"messages": messages,
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

// GetGroupMembers -> Returns all the members (users) of the group
func (handler *Handler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	if _, errResp := utils.CheckAuth(r.Header, handler.Paseto); errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "group_id")
	if groupId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "group id is missing")
		return
	}

	groupObjectId, errResp := utils.ToObjectId(groupId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"_id": groupObjectId,
	}

	projection := bson.M{
		"members": 1,
	}

	groupInstance, err := handler.Models.Group.Get(filter, projection)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "fetchMessages", err)
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

	fmt.Println(members)

	utils.WriteJSON(w, http.StatusOK, members)
}

func (handler *Handler) BanMemberFromGroup(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "group_id")
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

	filter := bson.M{
		"_id": groupObjectId,
	}

	projection := bson.M{
		"members":        1,
		"banned_members": 1,
		"owner_id":       1,
	}

	groupInstance, err := handler.Models.Group.Get(filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getGroup", "group with this id does not exist")
			return
		}

		utils.WriteError(w, http.StatusBadRequest, "getGroup", err.Error())
		return
	}

	if !slices.Contains(groupInstance.Admins, payload.UserId) {
		utils.WriteError(w, http.StatusBadRequest, "banFromGroup", "only group admins can ban someone")
		return
	}

	if slices.Contains(groupInstance.BannedMembers, targetUserObjectId) {
		utils.WriteError(w, http.StatusBadRequest, "banFromGroup", "this user is already banned from this group")
		return
	}

	if !slices.Contains(groupInstance.Members, targetUserObjectId) {
		utils.WriteError(w, http.StatusBadRequest, "banFromGroup", "this user is not a member of this group")
		return
	}

	newMembers := utils.DeleteElementFromSlice(groupInstance.Members, targetUserObjectId)
	newBans := append(groupInstance.BannedMembers, targetUserObjectId)

	updates := bson.M{
		"members":        newMembers,
		"banned_members": newBans,
	}

	if _, err := handler.Models.Group.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "updatingGroup", err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, "user banned from this group successfully")
}

func (handler *Handler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	groupId := chi.URLParam(r, "group_id")
	if groupId == "" {
		utils.WriteError(w, http.StatusBadRequest, "paramMissing", "group id is missing")
		return
	}

	groupObjectId, errResp := utils.ToObjectId(groupId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"_id": groupObjectId,
	}

	projection := bson.M{
		"members":  1,
		"owner_id": 1,
	}

	groupInstance, err := handler.Models.Group.Get(filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getGroup", "group with this id does not exist")
			return
		}

		utils.WriteError(w, http.StatusBadRequest, "getGroup", err.Error())
		return
	}

	if payload.UserId == groupInstance.OwnerId {
		utils.WriteError(w, http.StatusBadRequest, "leaveGroup", "group owner can`t leave it. You can Delete it")
		return
	}

	if !slices.Contains(groupInstance.Members, payload.UserId) {
		utils.WriteError(w, http.StatusBadRequest, "leaveGroup", "you are not a member of this group")
		return
	}

	newMembers := utils.DeleteElementFromSlice(groupInstance.Members, payload.UserId)

	updates := bson.M{
		"members": newMembers,
	}

	if _, err := handler.Models.Group.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "updatingGroup", err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, "you left the group successfully")
}

func checkUserApproval(groupId, userId primitive.ObjectID, handler *Handler) *utils.ErrorResponse {
	filter := bson.M{
		"group_id":     groupId,
		"requester_id": userId,
	}

	projection := bson.M{
		"status": 1,
	}

	approvalInstance, err := handler.Models.Approval.Get(filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &utils.ErrorResponse{Type: "userApprovalNotFound", Detail: "You have to submit an approval for this"}
		}

		return &utils.ErrorResponse{Type: "getUserApproval", Detail: err.Error()}
	}

	if approvalInstance.Status == "approved" {
		return nil
	}

	if approvalInstance.Status == "pending" {
		return &utils.ErrorResponse{Type: "userApprovalStatus", Detail: "Your approval is in pending status. Please be patient"}
	}

	if approvalInstance.Status == "rejected" {
		return &utils.ErrorResponse{Type: "userApprovalStatus", Detail: "Your approval has been rejected"}
	}

	return nil
}

func (handler *Handler) getIdByInviteLink(inviteLink string) (primitive.ObjectID, error) {
	filter := bson.M{
		"invite_link": inviteLink,
	}

	projection := bson.M{
		"_id": 1,
	}

	groupInstance, err := handler.Models.Group.Get(filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return primitive.NilObjectID, errors.New("group with this id does not exist")
		}

		return primitive.NilObjectID, errors.New(err.Error())
	}

	return groupInstance.Id, nil
}
