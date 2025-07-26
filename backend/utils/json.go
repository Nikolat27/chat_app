package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

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

func ParseJSON(reqBody io.ReadCloser, maxBytes int64, input any) error {
	body, err := io.ReadAll(io.LimitReader(reqBody, maxBytes))
	if err != nil {
		return err
	}

	return json.Unmarshal(body, input)
}
