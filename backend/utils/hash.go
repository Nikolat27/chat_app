package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(plainText []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(plainText, bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, nil
	}

	return hash, nil
}

func VerifyHash(hash, plainText string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainText)) == nil
}
