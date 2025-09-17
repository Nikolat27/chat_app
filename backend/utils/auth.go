package utils

import (
	"chat_app/paseto"
	"errors"
	"net/http"
)

var (
	NoAuthCookie   = errors.New("noCookieToken")
	CookieNotValid = errors.New("cookieNotValid")
)

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
