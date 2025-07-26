package paseto

import (
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Payload struct {
	ID        uuid.UUID          `json:"id"`
	UserId    primitive.ObjectID `json:"user_id"`
	Username  string             `json:"username"`
	CreatedAt time.Time          `json:"created_at"`
	ExpiryAt  time.Time          `json:"expiry_at"`
}

func NewPayload(userId primitive.ObjectID, username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		UserId:    userId,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiryAt:  time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiryAt) {
		return errors.New("token has expired")
	}

	return nil
}
