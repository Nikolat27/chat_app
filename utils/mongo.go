package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func ToObjectId(str string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(str)
}
