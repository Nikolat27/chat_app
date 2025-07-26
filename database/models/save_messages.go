package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type SaveMessageModel struct {
	collection *mongo.Collection
}

func NewSaveMessageModel(db *mongo.Database) *SaveMessageModel {
	return &SaveMessageModel{
		collection: db.Collection("save_messages"),
	}
}

type SaveMessage struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerId        primitive.ObjectID `json:"owner_id" bson:"owner_id"`
	Type           string             `json:"type" bson:"type"`
	Content        string             `json:"content" bson:"content"`
	ContentAddress string             `json:"content_address" bson:"content_address"`
	EditedAt       *time.Time         `json:"edited_at" bson:"edited_at"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}

func (save *SaveMessageModel) Create(ownerId primitive.ObjectID, contentType, content, contentAddress string) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newMessage = &SaveMessage{
		OwnerId:        ownerId,
		Content:        content,
		Type:           contentType,
		ContentAddress: contentAddress,
		CreatedAt:      time.Now(),
	}

	return save.collection.InsertOne(ctx, newMessage)
}

func (save *SaveMessageModel) Get(filter, projection bson.M) (*SaveMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	var msgInstance SaveMessage

	if err := save.collection.FindOne(ctx, filter, findOptions).Decode(&msgInstance); err != nil {
		return nil, err
	}

	return &msgInstance, nil
}

func (save *SaveMessageModel) GetAll(filter, projection bson.M, page, pageLimit int64) ([]SaveMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetProjection(projection)
	findOptions.SetSkip((page - 1) * pageLimit)
	findOptions.SetLimit(pageLimit)

	var messages []SaveMessage
	cursor, err := save.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (save *SaveMessageModel) Delete(filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return save.collection.DeleteOne(ctx, filter)
}

func (save *SaveMessageModel) Update(filter, updates bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
	}

	return save.collection.UpdateOne(ctx, filter, update)
}
