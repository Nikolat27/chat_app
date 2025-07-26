package handlers

import (
	"chat_app/utils"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func (handler *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.CheckAuth(r.Header, handler.Paseto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	file, header, err2 := r.FormFile("file")
	if err2 != nil {
		utils.WriteError(w, http.StatusBadRequest, "getFile", err2)
		return
	}
	defer file.Close()

	if err := os.MkdirAll("uploads", 0755); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "mkdirAll", err)
		return
	}

	fileName := rand.Text() + header.Filename
	path := filepath.Join("uploads", fileName)

	dst, err2 := os.Create(path)
	if err2 != nil {
		utils.WriteError(w, http.StatusBadRequest, "openDir", err2)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "copyFile", err)
		return
	}

	filter := bson.M{
		"_id": payload.UserId,
	}

	updates := bson.M{
		"avatar_url": fileName,
	}

	if _, err := handler.Models.User.Update(filter, updates); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "updateUser", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "user updated successfully")
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
