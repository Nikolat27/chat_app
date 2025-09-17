package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatModel struct {
	collection *mongo.Collection
}

// Chat -> This model is for 1 v 1 chat (direct)
type Chat struct {
	Id            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Participants  []primitive.ObjectID `json:"participants" bson:"participants"`
	LastMessageId primitive.ObjectID   `json:"last_message_id" bson:"last_message_id"`
	LastMessageAt time.Time            `json:"last_message_at" bson:"last_message_at"`
	CreatedAt     time.Time            `json:"created_at" bson:"created_at"`
}

func NewChatModel(db *mongo.Database) *ChatModel {
	return &ChatModel{
		collection: db.Collection("chats"),
	}
}

func (chat *ChatModel) Create(participants []primitive.ObjectID) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newChat = &Chat{
		Participants:  participants,
		LastMessageAt: time.Now(),
		CreatedAt:     time.Now(),
	}

	result, err := chat.collection.InsertOne(ctx, newChat)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (chat *ChatModel) Get(filter, projection bson.M) (*Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	var chatInstance Chat
	if err := chat.collection.FindOne(ctx, filter, findOptions).Decode(&chatInstance); err != nil {
		return nil, err
	}

	return &chatInstance, nil
}

func (chat *ChatModel) GetAll(filter bson.M, projection bson.M, page, pageLimit int64) ([]Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetProjection(projection)
	findOptions.SetSkip((page - 1) * pageLimit)
	findOptions.SetLimit(pageLimit)

	var chats []Chat
	cursor, err := chat.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &chats); err != nil {
		return nil, err
	}

	return chats, nil
}

func (chat *ChatModel) Delete(filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := chat.collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (chat *ChatModel) Update(filter bson.M, updates bson.M) (*mongo.UpdateResult, error) {
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
