package handlers

import "net/http"

func (handler *Handler) Register(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("register"))
}

func (handler *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login"))
}
