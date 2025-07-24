package handlers

import (
	"chat_app/utils"
	"net/http"
)

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, "register")
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusBadRequest, "login")
}
