package paseto

import (
	"crypto/sha256"
	"errors"
	"os"
	"time"

	"github.com/o1egl/paseto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Maker struct {
	paseto       *paseto.V2
	symmetricKey [32]byte
}

func New() (*Maker, error) {
	key, err := getSymmetricKey()
	if err != nil {
		return nil, err
	}

	maker := &Maker{
		paseto:       paseto.NewV2(),
		symmetricKey: key,
	}

	return maker, nil
}

func getSymmetricKey() ([32]byte, error) {
	symmetricKey := os.Getenv("PASETO_SYMMETRIC_KEY")
	if symmetricKey == "" {
		return [32]byte{}, errors.New("PASETO_SYMMETRIC_KEY env variable does not exist")
	}

	hash := sha256.Sum256([]byte(symmetricKey))

	return hash, nil
}

func (maker *Maker) CreateToken(userId primitive.ObjectID, username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(userId, username, duration)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmetricKey[:], payload, nil)
}

func (maker *Maker) VerifyToken(token string) (*Payload, error) {
	var payload Payload
	if err := maker.paseto.Decrypt(token, maker.symmetricKey[:], &payload, nil); err != nil {
		return nil, err
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return &payload, nil
}
