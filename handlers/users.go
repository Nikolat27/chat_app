package handlers

import (
	"chat_app/utils"
	"encoding/hex"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
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

func (handler *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
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
