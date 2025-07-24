package models

import "go.mongodb.org/mongo-driver/mongo"

type ChatModel struct {
	collection *mongo.Collection
}

func NewChatModel(db *mongo.Database) *ChatModel {
	return &ChatModel{
		collection: db.Collection("chats"),
	}
}

func (chat *ChatModel) Create() {
}
