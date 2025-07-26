package cipher

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/chacha20poly1305"
	"os"
)

type Cipher struct {
	Aead cipher.AEAD
}

func New() *Cipher {
	secretKey := os.Getenv("ENCRYPTION_SECRET_KEY")
	if secretKey == "" {
		panic("ENCRYPTION_SECRET_KEY env var is empty")
	}

	hashedKey := sha256.Sum256([]byte(secretKey))

	aead, err := chacha20poly1305.New(hashedKey[:])
	if err != nil {
		panic(err)
	}

	return &Cipher{
		Aead: aead,
	}
}

func (cipher *Cipher) Encrypt(plainText []byte) ([]byte, error) {
	nonce, err := generateNonce()
	if err != nil {
		return nil, err
	}

	cipherText := cipher.Aead.Seal(nil, nonce, plainText, nil)

	return append(nonce, cipherText...), nil
}

func (cipher *Cipher) Decrypt(ciphered []byte) ([]byte, error) {
	if len(ciphered) < 12 {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := ciphered[:12]
	cipherText := ciphered[12:]

	return cipher.Aead.Open(nil, nonce, cipherText, nil)
}

func generateNonce() ([]byte, error) {
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}
