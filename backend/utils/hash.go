package utils

import (
	"crypto/rand"
	"crypto/sha256"
)

const Size = 32

func Hash(plainText, salt []byte) [Size]byte {
	newHash := append(plainText, salt...)
	return sha256.Sum256(newHash)
}

func VerifyHash(hash, salt, plainText []byte) bool {
	combinedSlice := append(plainText, salt...)
	newHash := sha256.Sum256(combinedSlice)

	return string(newHash[:]) == string(hash)
}

func GenerateSalt(size int64) ([]byte, error) {
	newBytes := make([]byte, size)
	if _, err := rand.Read(newBytes); err != nil {
		return nil, err
	}

	return newBytes, nil
}
