package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ApprovalModel struct {
	collection *mongo.Collection
}

func NewApprovalModel(db *mongo.Database) *ApprovalModel {
	return &ApprovalModel{
		collection: db.Collection("approvals"),
	}
}

type Approval struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	GroupId      primitive.ObjectID `json:"group_id" bson:"group_id"`
	GroupOwnerId primitive.ObjectID `json:"group_owner_id" bson:"group_owner_id"`
	RequesterId  primitive.ObjectID `json:"requester_id" bson:"requester_id"`
	Status       string             `json:"status" bson:"status"` // pending, rejected, approved
	Reason       string             `json:"reason" bson:"reason"`
	ReviewedAt   *time.Time         `json:"reviewed_at" bson:"reviewed_at"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

func (approval *ApprovalModel) Create(groupId, groupOwnerId, requesterId primitive.ObjectID, reason string) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newApproval := &Approval{
		GroupId:      groupId,
		GroupOwnerId: groupOwnerId,
		RequesterId:  requesterId,
		Status:       "pending",
		Reason:       reason,
		CreatedAt:    time.Now(),
	}

	return approval.collection.InsertOne(ctx, newApproval)
}

func (approval *ApprovalModel) Update(filter, updates bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": updates,
	}

	return approval.collection.UpdateOne(ctx, filter, update)
}

func (approval *ApprovalModel) Get(filter, projection bson.M) (*Approval, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.FindOne()
	findOptions.SetProjection(projection)

	var approvalInstance Approval
	if err := approval.collection.FindOne(ctx, filter, findOptions).Decode(&approvalInstance); err != nil {
		return nil, err
	}

	return &approvalInstance, nil
}

func (approval *ApprovalModel) GetAll(filter, projection bson.M, page, pageLimit int64) ([]Approval, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetProjection(projection)
	findOptions.SetSkip((page - 1) * pageLimit)
	findOptions.SetLimit(pageLimit)

	var approvals []Approval
	cursor, err := approval.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &approvals); err != nil {
		return nil, err
	}

	return approvals, nil
}

func (approval *ApprovalModel) Delete(filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return approval.collection.DeleteOne(ctx, filter)
}

func (approval *ApprovalModel) DeleteAll(filter bson.M) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return approval.collection.DeleteMany(ctx, filter)
}
