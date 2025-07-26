package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type GroupModel struct {
	collection *mongo.Collection
}

func NewGroupModel(db *mongo.Database) *GroupModel {
	return &GroupModel{
		collection: db.Collection("groups"),
	}
}

type Group struct {
	Id              primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerId         primitive.ObjectID   `json:"owner_id" bson:"owner_id"`
	Users           []primitive.ObjectID `json:"users" bson:"users"`
	Name            string               `json:"name" bson:"name"`
	Description     string               `json:"description" bson:"description"`
	AvatarUrl       string               `json:"avatar_url" bson:"avatar_url"`
	InviteLink      string               `json:"invite_link" bson:"invite_link"`
	PinnedMessageId primitive.ObjectID   `json:"pinned_message_id" bson:"pinned_message_id"`
	LastMessageId   primitive.ObjectID   `json:"last_message_id" bson:"last_message_id"`
	LastMessageAt   time.Time            `json:"last_message_at" bson:"last_message_at"`
	CreatedAt       time.Time            `json:"created_at" bson:"created_at"`
}

func (group *GroupModel) Create(ownerId primitive.ObjectID, name, description,
	avatarUrl, inviteLink string, users []primitive.ObjectID) (*mongo.InsertOneResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newGroup := &Group{
		OwnerId:     ownerId,
		Name:        name,
		Description: description,
		AvatarUrl:   avatarUrl,
		InviteLink:  inviteLink,
		Users:       users,
		CreatedAt:   time.Now(),
	}

	return group.collection.InsertOne(ctx, newGroup)
}

func (group *GroupModel) Get(filter, projection bson.M) (*Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	var groupInstance Group
	if err := group.collection.FindOne(ctx, filter, findOptions).Decode(&groupInstance); err != nil {
		return nil, err
	}

	return &groupInstance, nil
}

func (group *GroupModel) Update(filter, updates bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
	}

	return group.collection.UpdateOne(ctx, filter, update)
}

func (group *GroupModel) Delete(filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	return group.collection.DeleteOne(ctx, filter)
}
