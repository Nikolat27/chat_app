package handlers

import (
	"chat_app/utils"
	"encoding/hex"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		RawPassword string `json:"password"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err)
		return
	}

	filter := bson.M{
		"username": input.Username,
		"email":    input.Email,
	}

	projection := bson.M{
		"_id": 1,
	}

	_, err := handler.Models.User.Get(filter, projection)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		utils.WriteError(w, http.StatusBadRequest, "getUserInstance", "either this username or email exists")
		return
	}

	salt, err := utils.GenerateSalt(14)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "generatingSalt", err)
		return
	}

	hashedPassword := utils.Hash(input.RawPassword, salt)

	encodedHashedPassword := hex.EncodeToString([]byte(hashedPassword))
	encodedSalt := hex.EncodeToString([]byte(salt))

	if _, err := handler.Models.User.Create(input.Username, input.Email, encodedHashedPassword, encodedSalt); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "creatingUserInstance", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "user registered successfully")
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusBadRequest, "login")
}
