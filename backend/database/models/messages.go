package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MessageModel struct {
	collection *mongo.Collection
}

func NewMessageModel(db *mongo.Database) *MessageModel {
	return &MessageModel{
		collection: db.Collection("messages"),
	}
}

type Message struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ChatId     primitive.ObjectID `json:"chat_id" bson:"chat_id"`
	GroupId    primitive.ObjectID `json:"group_id" bson:"group_id"`
	SenderId   primitive.ObjectID `json:"sender_id" bson:"sender_id"`
	ReceiverId primitive.ObjectID `json:"receiver_id" bson:"receiver_id"`
	// text or image
	Type    string `json:"type" bson:"type"`
	Content string `json:"content" bson:"content"`
	// used for image addresses
	ContentAddress     string     `json:"content_address" bson:"content_address"`
	IsSecret           bool       `json:"is_secret" bson:"is_secret"`
	IsDeletedForSender bool       `json:"is_deleted_for_sender" bson:"is_deleted_for_sender"`
	EditedAt           *time.Time `json:"edited_at" bson:"edited_at"`
	CreatedAt          time.Time  `json:"created_at" bson:"created_at"`
}

func (message *MessageModel) Create(chatId, groupId, senderId, receiverId primitive.ObjectID, contentType, contentAddress,
	content string, isSecret bool) (*mongo.InsertOneResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newUser = &Message{
		ChatId:         chatId,
		GroupId:        groupId,
		SenderId:       senderId,
		ReceiverId:     receiverId,
		Content:        content,
		Type:           contentType,
		ContentAddress: contentAddress,
		IsSecret:       isSecret,
		CreatedAt:      time.Now(),
	}

	return message.collection.InsertOne(ctx, newUser)
}

func (message *MessageModel) GetAll(filter, projection bson.M, page, pageLimit int64) ([]Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetProjection(projection)
	findOptions.SetSkip((page - 1) * pageLimit)
	findOptions.SetLimit(pageLimit)
	findOptions.SetSort(bson.M{
		"created_at": 1,
	})

	var messages []Message
	cursor, err := message.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (message *MessageModel) Get(filter, projection bson.M) (*Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetProjection(projection)

	var messageInstance Message
	if err := message.collection.FindOne(ctx, filter).Decode(&messageInstance); err != nil {
		return nil, err
	}

	return &messageInstance, nil
}

func (message *MessageModel) Delete(filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return message.collection.DeleteOne(ctx, filter)
}

func (message *MessageModel) DeleteAll(filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return message.collection.DeleteMany(ctx, filter)
}

func (message *MessageModel) Update(filter, updates bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
	}

	return message.collection.UpdateOne(ctx, filter, update)
}
