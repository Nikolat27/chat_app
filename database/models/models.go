package models

import "go.mongodb.org/mongo-driver/mongo"

type Models struct {
	User    *UserModel
	Chat    *ChatModel
	Message *MessageModel
}

func New(db *mongo.Database) *Models {
	return &Models{
		User:    NewUserModel(db),
		Chat:    NewChatModel(db),
		Message: NewMessageModel(db),
	}
}
