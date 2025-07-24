package models

import "go.mongodb.org/mongo-driver/mongo"

type UserModel struct {
	collection *mongo.Collection
}

func NewUserModel(db *mongo.Database) *UserModel {
	return &UserModel{
		collection: db.Collection("users"),
	}
}

func (user *UserModel) Create() {
}
