package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type SecretChatModel struct {
	collection *mongo.Collection
}

// SecretChat -> This model is for e2ee 1 v 1 chat (direct)
type SecretChat struct {
	Id                         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	User1                      primitive.ObjectID `json:"user_1" bson:"user_1"`
	User2                      primitive.ObjectID `json:"user_2" bson:"user_2"`
	User1EncryptedSymmetricKey string             `json:"user_1_encrypted_symmetric_key" bson:"user_1_encrypted_symmetric_key"`
	User2EncryptedSymmetricKey string             `json:"user_2_encrypted_symmetric_key" bson:"user_2_encrypted_symmetric_key"`
	User2Accepted              bool               `json:"user_2_accepted" bson:"user_2_accepted"`
	ExpireAt                   *time.Time         `json:"expire_at" bson:"expire_at"`
	CreatedAt                  time.Time          `json:"created_at" bson:"created_at"`
}

func NewSecretChatModel(db *mongo.Database) *SecretChatModel {
	return &SecretChatModel{
		collection: db.Collection("secret_chats"),
	}
}

func (chat *SecretChatModel) Create(user1, user2 primitive.ObjectID) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newChat := &SecretChat{
		User1:     user1,
		User2:     user2,
		CreatedAt: time.Now(),
	}

	result, err := chat.collection.InsertOne(ctx, newChat)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (chat *SecretChatModel) Get(filter, projection bson.M) (*SecretChat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.FindOne().SetProjection(projection)

	var chatInstance SecretChat
	if err := chat.collection.FindOne(ctx, filter, findOptions).Decode(&chatInstance); err != nil {
		return nil, err
	}

	return &chatInstance, nil
}

func (chat *SecretChatModel) GetAll(filter, projection bson.M, page, pageLimit int64) ([]SecretChat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find().
		SetProjection(projection).
		SetSkip((page - 1) * pageLimit).
		SetLimit(pageLimit)

	var chats []SecretChat
	cursor, err := chat.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &chats); err != nil {
		return nil, err
	}

	return chats, nil
}

func (chat *SecretChatModel) Delete(filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := chat.collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (chat *SecretChatModel) Update(filter bson.M, updates bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
	}

	result, err := chat.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}
