package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	Admins          []primitive.ObjectID `json:"admins" bson:"admins"`
	Members         []primitive.ObjectID `json:"members" bson:"members"`
	BannedMembers   []primitive.ObjectID `json:"banned_members" bson:"banned_members"`
	Name            string               `json:"name" bson:"name"`
	Description     string               `json:"description" bson:"description"`
	AvatarUrl       string               `json:"avatar_url" bson:"avatar_url"`
	Type            string               `json:"type" bson:"type"` // public or private (private needs apporval)
	InviteLink      string               `json:"invite_link" bson:"invite_link"`
	PinnedMessageId primitive.ObjectID   `json:"pinned_message_id" bson:"pinned_message_id"`
	LastMessageId   primitive.ObjectID   `json:"last_message_id" bson:"last_message_id"`
	IsSecret        bool                 `json:"is_secret" bson:"is_secret"`
	LastMessageAt   time.Time            `json:"last_message_at" bson:"last_message_at"`
	CreatedAt       time.Time            `json:"created_at" bson:"created_at"`
}

func (group *GroupModel) Create(ownerId primitive.ObjectID, name, description, avatarUrl, groupType,
	inviteLink string, members, admins []primitive.ObjectID, isSecret bool) (*mongo.InsertOneResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newGroup := &Group{
		OwnerId:     ownerId,
		Name:        name,
		Description: description,
		AvatarUrl:   avatarUrl,
		Type:        groupType,
		InviteLink:  inviteLink,
		Members:     members,
		Admins:      admins,
		IsSecret:    isSecret,
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

func (group *GroupModel) GetAll(filter, projection bson.M, page, pageLimit int64) ([]Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetProjection(projection)
	findOptions.SetSkip((page - 1) * pageLimit)
	findOptions.SetLimit(pageLimit)

	var groups []Group
	cursor, err := group.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}

	return groups, nil
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
