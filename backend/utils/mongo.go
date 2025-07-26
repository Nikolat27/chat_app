package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func ToObjectId(str string) (primitive.ObjectID, *ErrorResponse) {
	objectId, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return primitive.NilObjectID, &ErrorResponse{
			Type:   "checkAuth",
			Detail: err,
		}
	}

	return objectId, nil
}
