package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UserModel struct {
	collection *mongo.Collection
}

type User struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username       string             `json:"username" bson:"username"`
	Email          string             `json:"email" bson:"email"`
	HashedPassword string             `json:"hashed_password" bson:"hashed_password"`
	Salt           string             `json:"salt" bson:"salt"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}

func NewUserModel(db *mongo.Database) *UserModel {
	collection := db.Collection("users")

	_, err := collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})

	if err != nil {
		panic(fmt.Errorf("ERROR creating index on users: %s", err))
	}

	return &UserModel{
		collection: collection,
	}
}

func (user *UserModel) Create(username, email, hashedPassword, salt string) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newUser = &User{
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
		Salt:           salt,
		CreatedAt:      time.Now(),
	}

	result, err := user.collection.InsertOne(ctx, newUser)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// Get -> Returns One
func (user *UserModel) Get(filter bson.M, projection bson.M) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	var userInstance User
	if err := user.collection.FindOne(ctx, filter, findOptions).Decode(&userInstance); err != nil {
		return nil, err
	}

	return &userInstance, nil
}
