package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment is the data structure that we will save and receive.
type Like struct {
	ID          primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Like_ID     interface{}            `json:"like_id,omitempty" bson:"like_id,omitempty" validate:"unique"`
	Description string                 `json:"description,omitempty" bson:"description,omitempty" validate:"required"`
	Person_Info Person                 `json:"person_info,omitempty" bson:"person_info,omitempty"`
	Post_Info   Blog                   `json:"post_id,omitempty" bson:"post_id,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"` // data is a optional fields that can hold anything in key:value format.
}
