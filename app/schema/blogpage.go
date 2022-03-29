package schema

import (
	//"blog/app/schema"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Blog is the data structure that we will save and receive.
//now := time.Now()
type Blog struct {
	ID          primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string                 `json:"title,omitempty" bson:"title,omitempty"`
	Description string                 `json:"description,omitempty" bson:"description,omitempty"`
	User_ID     primitive.ObjectID     `json:"user_id,omitempty" bson:"user_id,omitempty" validate:"required"`
	BlogImg     string                 `json:"blog_img,omitempty" bson:"blog_img,omitempty"`
	Tag         string                 `json:"tag,omitempty" bson:"tag,omitempty"`
	Time        string                 `json:"time,omitempty" bson:"time,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"` // data is a optional fields that can hold anything in key:value format.
}
