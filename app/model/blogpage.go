package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Person is the data structure that we will save and receive.
type Blog struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Author      string             `json:"author,omitempty" bson:"author,omitempty"`
	//Username    string                 `json:"username,omitempty" bson:"username,omitempty"`
	//Email       string                 `json:"email,omitempty" bson:"email,omitempty"`
	Data map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"` // data is a optional fields that can hold anything in key:value format.
}
