package handlers

import (
	"chat_app/utils"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username    string `json:"username"`
		RawPassword string `json:"password"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parseJson", err)
		return
	}

	filter := bson.M{
		"username": input.Username,
	}

	projection := bson.M{
		"_id": 1,
	}

	if _, err := handler.Models.User.Get(filter, projection); err == nil {
		utils.WriteError(w, http.StatusBadRequest, "usernameEmailTaken", "user with this username exists already")
		return
	}

	if err := handler.CreateUser(input.Username, []byte(input.RawPassword)); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Type, err.Detail)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "user registered successfully")
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username    string `json:"username"`
		RawPassword string `json:"password"`
	}

	if err := utils.ParseJSON(r.Body, 1_000, &input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "parsingJson", err)
		return
	}

	filter := bson.M{
		"username": input.Username,
	}

	projection := bson.M{
		"_id":             1,
		"hashed_password": 1,
		"avatar_url":      1,
	}

	user, err := handler.Models.User.Get(filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			utils.WriteError(w, http.StatusBadRequest, "userDoesNotExist", "user with this username does not exist")
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, "userGetError", "failed to fetch user")
		return
	}

	if !utils.VerifyHash(user.HashedPassword, input.RawPassword) {
		utils.WriteError(w, http.StatusBadRequest, "passwordValidation", "password is invalid")
		return
	}

	token, err := handler.Paseto.CreateToken(user.Id, input.Username, 12*time.Hour)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "createPasetoToken", "failed to create paseto token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_cookie",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // temp (for http)
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600 * 12, // 12 hours
	})

	var response = map[string]string{
		"username":   input.Username,
		"user_id":    user.Id.Hex(),
		"avatar_url": user.AvatarUrl,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (handler *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_cookie",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // delete it
	})

	utils.WriteJSON(w, http.StatusOK, "Logged out successfully")
}

func (handler *Handler) AuthCheck(w http.ResponseWriter, r *http.Request) {
	payload, errResp := utils.CheckAuth(r, handler.Paseto)
	if errResp != nil {
		utils.WriteError(w, http.StatusUnauthorized, errResp.Type, errResp.Detail)
		return
	}

	var response = map[string]string{
		"username": payload.Username,
		"user_id":  payload.UserId.Hex(),
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
