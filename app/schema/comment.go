package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment is the data structure that we will save and receive.
type Comment struct {
	ID          primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Description string                 `json:"description,omitempty" bson:"description,omitempty" validate:"required"`
	User_ID     primitive.ObjectID     `bson:"user_id,omitempty" validate:"required"`
	Post_ID     primitive.ObjectID     `bson:"post_id,omitempty" validate:"required"`
	Data        map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"` // data is a optional fields that can hold anything in key:value format.
}
