package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type SecretGroupMessageModel struct {
	collection *mongo.Collection
}

func NewSecretGroupMessageModel(db *mongo.Database) *SecretGroupMessageModel {
	return &SecretGroupMessageModel{
		collection: db.Collection("secret_group_messages"),
	}
}

type SecretGroupMessage struct {
	Id                 primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	GroupId            primitive.ObjectID `json:"group_id" bson:"group_id"`
	SenderId           primitive.ObjectID `json:"sender_id" bson:"sender_id"`
	EncryptedContent   string             `json:"encrypted_content" bson:"encrypted_content"` // base64
	IsDeletedForSender bool               `json:"is_deleted_for_sender" bson:"is_deleted_for_sender"`
	EditedAt           *time.Time         `json:"edited_at,omitempty" bson:"edited_at,omitempty"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
}

func (message *SecretGroupMessageModel) Create(groupId, senderId primitive.ObjectID,
	content string) (*mongo.InsertOneResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newMessage := &SecretGroupMessage{
		GroupId:          groupId,
		SenderId:         senderId,
		EncryptedContent: content,
		CreatedAt:        time.Now(),
	}

	return message.collection.InsertOne(ctx, newMessage)
}

func (message *SecretGroupMessageModel) GetAll(filter, projection bson.M, page, pageLimit int64) ([]SecretGroupMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetProjection(projection)
	findOptions.SetSkip((page - 1) * pageLimit)
	findOptions.SetLimit(pageLimit)
	findOptions.SetSort(bson.M{
		"created_at": 1,
	})

	var messages []SecretGroupMessage
	cursor, err := message.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
