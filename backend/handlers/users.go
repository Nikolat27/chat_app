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

func (handler *Handler) CreateUser(username string, rawPassword []byte) *utils.ErrorResponse {
	salt, err := utils.GenerateSalt(14)
	if err != nil {
		return &utils.ErrorResponse{
			Type:   "generatingSalt",
			Detail: err,
		}
	}

	hashedPassword := utils.Hash(rawPassword, salt)

	encodedHashedPassword := hex.EncodeToString(hashedPassword[:])
	encodedSalt := hex.EncodeToString(salt)

	if _, err := handler.Models.User.Create(username, encodedHashedPassword, encodedSalt); err != nil {
		return &utils.ErrorResponse{
			Type:   "createUser",
			Detail: "failed to create user",
		}
	}

	return nil
}

func (handler *Handler) SearchUser(w http.ResponseWriter, r *http.Request) {
	if _, errResp := utils.CheckAuth(r.Header, handler.Paseto); errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		utils.WriteError(w, http.StatusBadRequest, "missingQuery", "q query is missing")
		return
	}

	filter := bson.M{
		"username": query,
	}

	userInstance, err := handler.Models.User.Get(filter, bson.M{})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		utils.WriteError(w, http.StatusBadRequest, "getUser", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, userInstance)
}

func (handler *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	if _, errResp := utils.CheckAuth(r.Header, handler.Paseto); errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		utils.WriteError(w, http.StatusBadRequest, "getUrlParam", "user id is missing")
		return
	}

	userObjectId, errResp := utils.ToObjectId(userId)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"_id": userObjectId,
	}

	projection := bson.M{
		"avatar_url": 1,
		"username":   1,
	}

	userInstance, err := handler.Models.User.Get(filter, projection)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		utils.WriteError(w, http.StatusBadRequest, "getUser", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, userInstance)
}

func (handler *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Type, err.Detail)
		return
	}

	allowedFormats := []string{".png", ".jpeg", ".jpg", ".webp"}

	avatarAddress, err := utils.UploadFile(r, "file", allowedFormats)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	filter := bson.M{
		"_id": payload.UserId,
	}

	updates := bson.M{
		"avatar_url": avatarAddress,
	}

	if _, err := handler.Models.User.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "updateUser", err)
		return
	}

	response := map[string]string{
		"avatar_url": avatarAddress,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (handler *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, err.Type, err.Detail)
		return
	}

	filter := bson.M{
		"_id": payload.UserId,
	}

	projection := bson.M{
		"_id": 1,
	}

	if _, err := handler.Models.User.Get(filter, projection); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "getUser", "user with this id does not exist")
			return
		}

		utils.WriteError(w, http.StatusBadRequest, "getUser", "failed to get user")
		return
	}

	if _, err := handler.Models.User.Delete(filter); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "deleteUser", "failed to delete user")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "user deleted successfully")
}

// GetUserChats -> Returns the chats themselves
func (handler *Handler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"participants": bson.M{
			"$in": []primitive.ObjectID{payload.UserId},
		},
	}

	chats, err := handler.Models.Chat.GetAll(filter, bson.M{}, 1, 10)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		utils.WriteError(w, http.StatusBadRequest, "getChats", err)
		return
	}

	avatarUrls := make(map[string]string)
	usernames := make(map[string]string)
	for _, chat := range chats {
		otherUserId := getOtherUserId(chat.Participants, payload.UserId)
		url, _ := getUserAvatarUrl(otherUserId, handler)
		username, _ := getUserUsername(otherUserId, handler)

		avatarUrls[chat.Id.Hex()] = url
		usernames[chat.Id.Hex()] = username
	}

	response := map[string]any{
		"chats":       chats,
		"avatar_urls": avatarUrls,
		"usernames":   usernames,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetUserSecretChats -> Returns the secret chats themselves
func (handler *Handler) GetUserSecretChats(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"$or": []bson.M{
			{"user_1": payload.UserId},
			{"user_2": payload.UserId},
		},
	}

	chats, err := handler.Models.SecretChat.GetAll(filter, bson.M{}, 1, 10)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		utils.WriteError(w, http.StatusBadRequest, "getChats", err)
		return
	}

	usernames := make(map[string]string)
	for _, chat := range chats {
		otherUserId := getOtherUserId([]primitive.ObjectID{chat.User1, chat.User2}, payload.UserId)
		username, _ := getUserUsername(otherUserId, handler)

		usernames[chat.Id.Hex()] = username
	}

	response := map[string]any{
		"secret_chats":     chats,
		"secret_usernames": usernames,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetUserGroups -> Returns the chats themselves
func (handler *Handler) GetUserGroups(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"members": bson.M{
			"$in": []primitive.ObjectID{payload.UserId},
		},
	}

	page, pageLimit, errResp := utils.ParsePageAndLimitQueryParams(r.URL)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	groups, err := handler.Models.Group.GetAll(filter, bson.M{}, page, pageLimit)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		utils.WriteError(w, http.StatusBadRequest, "getGroups", err.Error())
		return
	}

	response := map[string]any{
		"groups": groups,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetUserSecretGroups -> Returns the secret chats themselves
func (handler *Handler) GetUserSecretGroups(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r.Header, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	filter := bson.M{
		"members": bson.M{
			"$in": []primitive.ObjectID{payload.UserId},
		},
	}

	page, pageLimit, errResp := utils.ParsePageAndLimitQueryParams(r.URL)
	if errResp != nil {
		utils.WriteError(w, http.StatusBadRequest, errResp.Type, errResp.Detail)
		return
	}

	groups, err := handler.Models.SecretGroup.GetAll(filter, bson.M{}, page, pageLimit)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		utils.WriteError(w, http.StatusBadRequest, "getGroups", err.Error())
		return
	}

	response := map[string]any{
		"groups": groups,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func getOtherUserId(participants []primitive.ObjectID, userId primitive.ObjectID) primitive.ObjectID {
	for _, participant := range participants {
		if participant != userId {
			return participant
		}
	}

	return userId
}

func getUserAvatarUrl(id primitive.ObjectID, handler *Handler) (string, error) {
	filter := bson.M{
		"_id": id,
	}

	projection := bson.M{
		"avatar_url": 1,
	}

	user, err := handler.Models.User.Get(filter, projection)
	if err != nil {
		return "", err
	}

	return user.AvatarUrl, nil
}

func getUserUsername(id primitive.ObjectID, handler *Handler) (string, error) {
	filter := bson.M{
		"_id": id,
	}

	projection := bson.M{
		"username": 1,
	}

	user, err := handler.Models.User.Get(filter, projection)
	if err != nil {
		return "", err
	}

	return user.Username, nil
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
