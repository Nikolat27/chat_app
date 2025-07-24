package models

import "go.mongodb.org/mongo-driver/mongo"

type Models struct {
	User *UserModel
	Chat *ChatModel
}

func New(db *mongo.Database) *Models {
	return &Models{
		User: NewUserModel(db),
		Chat: NewChatModel(db),
	}
}
