package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Person is the data structure that we will save and receive.
type Blog struct {
	ID          primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Author_Id   primitive.ObjectID     `json: "author_id,omitempty" bson: "author_id,omitempty"`
	Title       string                 `json:"title,omitempty" bson:"title,omitempty"`
	Description string                 `json:"description,omitempty" bson:"description,omitempty"`
	Author_info Person                 `json:"author_info,omitempty" bson:"author_info,omitempty"`
	Username    string                 `json:"username,omitempty" bson:"username,omitempty"`
	Email       string                 `json:"email,omitempty" bson:"email,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"` // data is a optional fields that can hold anything in key:value format.
}
