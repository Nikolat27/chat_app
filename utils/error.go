package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Type   string `json:"type,omitempty"`
	Detail any    `json:"detail,omitempty"`
}

func WriteError(w http.ResponseWriter, status int, errType any, errDetail any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var errorType error

	switch value := errType.(type) {
	case string:
		errorType = errors.New(value)
	case error:
		errorType = value
	default:
		errorType = fmt.Errorf("unexepected error type :%v", errorType)
	}

	var resp = ErrorResponse{
		Type:   errorType.Error(),
		Detail: errDetail,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
