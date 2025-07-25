package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ChatModel struct {
	collection *mongo.Collection
}

type Chat struct {
	Id            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Participants  []primitive.ObjectID `json:"participants" bson:"participants"`
	LastMessageId primitive.ObjectID   `json:"last_message_id" bson:"last_message_id"`
	UpdatedAt     time.Time            `json:"updated_at" bson:"updated_at"`
	CreatedAt     time.Time            `json:"created_at" bson:"created_at"`
}

func NewChatModel(db *mongo.Database) *ChatModel {
	return &ChatModel{
		collection: db.Collection("chats"),
	}
}

func (chat *ChatModel) Create() {
}
