package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type SecretGroupModel struct {
	collection *mongo.Collection
}

func NewSecretGroupModel(db *mongo.Database) *SecretGroupModel {
	return &SecretGroupModel{
		collection: db.Collection("secret_groups"),
	}
}

type SecretGroup struct {
	Id              primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerId         primitive.ObjectID   `json:"owner_id" bson:"owner_id"`
	Admins          []primitive.ObjectID `json:"admins" bson:"admins"`
	Members         []primitive.ObjectID `json:"members" bson:"members"`
	BannedMembers   []primitive.ObjectID `json:"banned_members" bson:"banned_members"`
	UserPublicKeys  map[string]string    `json:"user_public_keys" bson:"user_public_keys"` // userId -> publicKey
	MemberJoinTimes map[string]time.Time `json:"join_times" bson:"join_times"`             // userId -> joinedAt
	Name            string               `json:"name" bson:"name"`
	Description     string               `json:"description" bson:"description"`
	Type            string               `json:"type" bson:"type"` // public or private
	InviteLink      string               `json:"invite_link" bson:"invite_link"`
	PinnedMessageId primitive.ObjectID   `json:"pinned_message_id" bson:"pinned_message_id"`
	LastMessageId   primitive.ObjectID   `json:"last_message_id" bson:"last_message_id"`
	LastMessageAt   time.Time            `json:"last_message_at" bson:"last_message_at"`
	CreatedAt       time.Time            `json:"created_at" bson:"created_at"`
}

func (model *SecretGroupModel) Create(ownerId primitive.ObjectID, name, description, groupType, inviteLink string,
	members, admins []primitive.ObjectID, userJoins map[string]time.Time, userPublicKeys map[string]string) (*mongo.InsertOneResult, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newGroup := &SecretGroup{
		OwnerId:         ownerId,
		Name:            name,
		Description:     description,
		Type:            groupType,
		InviteLink:      inviteLink,
		Members:         members,
		Admins:          admins,
		MemberJoinTimes: userJoins,
		UserPublicKeys:  userPublicKeys,
		CreatedAt:       time.Now(),
	}

	return model.collection.InsertOne(ctx, newGroup)
}

func (model *SecretGroupModel) Get(filter, projection bson.M) (*SecretGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	var groupInstance SecretGroup
	if err := model.collection.FindOne(ctx, filter, findOptions).Decode(&groupInstance); err != nil {
		return nil, err
	}

	return &groupInstance, nil
}

func (model *SecretGroupModel) GetAll(filter, projection bson.M, page, pageLimit int64) ([]SecretGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetProjection(projection)
	findOptions.SetSkip((page - 1) * pageLimit)
	findOptions.SetLimit(pageLimit)

	var groups []SecretGroup
	cursor, err := model.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

func (model *SecretGroupModel) Update(filter, updates bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
	}

	return model.collection.UpdateOne(ctx, filter, update)
}

func (model *SecretGroupModel) Delete(filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return model.collection.DeleteOne(ctx, filter)
}
