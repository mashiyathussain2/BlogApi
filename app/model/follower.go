package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment is the data structure that we will save and receive.
type Follower struct {
	ID           primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	User_ID      primitive.ObjectID     `json:"user_id,omitempty" bson:"user_id,omitempty" validate:"required"`
	Following_ID primitive.ObjectID     `json:"following_id,omitempty" bson:"following_id,omitempty" validate:"required"`
	Data         map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"` // data is a optional fields that can hold anything in key:value format.
}
