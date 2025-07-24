package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Type   string `json:"type"`
	Detail any    `json:"detail,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encodedData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(encodedData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, status int, errorType string, errorDetail any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var response = ErrorResponse{
		Type:   errorType,
		Detail: errorDetail,
	}

	encodedData, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(encodedData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
