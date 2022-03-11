package schema

import (
	//"blog/app/schema"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog is the data structure that we will save and receive.
type Blog struct {
	ID          primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string                 `json:"title,omitempty" bson:"title,omitempty"`
	Description string                 `json:"description,omitempty" bson:"description,omitempty"`
	User_ID     primitive.ObjectID     `json:"user_id,omitempty" bson:"user_id,omitempty" validate:"required"`
	Data        map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"` // data is a optional fields that can hold anything in key:value format.
}