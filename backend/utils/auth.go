package utils

import (
	"chat_app/paseto"
	"errors"
	"net/http"
)

var (
	// NoAuthToken   = errors.New("noAuthToken")
	// TokenNotValid = errors.New("tokenNotValid")

	NoAuthCookie   = errors.New("noCookieToken")
	CookieNotValid = errors.New("cookieNotValid")
)

// func CheckAuth(header http.Header, paseto *paseto.Maker) (*paseto.Payload, *ErrorResponse) {
// 	authToken := header.Get("Authorization")
// 	if authToken == "" {
// 		return nil, &ErrorResponse{Type: NoAuthToken.Error(), Detail: "Authorization token is missing"}
// 	}

// 	payload, err := paseto.VerifyToken(authToken)
// 	if err != nil {
// 		return nil, &ErrorResponse{Type: TokenNotValid.Error(), Detail: "can`t verify the token"}
// 	}

// 	return payload, nil
// }

func CheckAuth(r *http.Request, paseto *paseto.Maker) (*paseto.Payload, *ErrorResponse) {
	cookie, err := r.Cookie("auth_cookie")
	if err != nil {
		return nil, &ErrorResponse{Type: NoAuthCookie.Error(), Detail: "Auth cookie is missing"}
	}

	payload, err := paseto.VerifyToken(cookie.Value)
	if err != nil {
		return nil, &ErrorResponse{Type: CookieNotValid.Error(), Detail: "Can't verify auth token"}
	}

	return payload, nil
}
