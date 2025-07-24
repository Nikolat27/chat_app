package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
)

func Hash(plainText, salt string) string {
	hash := sha256.Sum256([]byte(plainText + salt))
	return string(hash[:])
}

func VerifyHash(hash []byte, plainText, salt string) bool {
	newHash := sha256.Sum256([]byte(plainText + salt))

	return bytes.Equal(hash, newHash[:])
}

func GenerateSalt(size int64) (string, error) {
	newBytes := make([]byte, size)
	if _, err := rand.Read(newBytes); err != nil {
		return "", err
	}

	return string(newBytes), nil
}
