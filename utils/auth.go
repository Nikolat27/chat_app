package utils

import (
	"chat_app/paseto"
	"errors"
	"net/http"
)

var (
	NoAuthToken   = errors.New("noAuthToken")
	TokenNotValid = errors.New("tokenNotValid")
)

func CheckAuth(header http.Header, paseto *paseto.Maker) (*paseto.Payload, *ErrorResponse) {
	authToken := header.Get("Authorization")
	if authToken == "" {
		return nil, &ErrorResponse{Type: NoAuthToken, Detail: "Authorization token is missing"}
	}

	payload, err := paseto.VerifyToken(authToken)
	if err != nil {
		return nil, &ErrorResponse{Type: TokenNotValid, Detail: err}
	}

	return payload, nil
}
