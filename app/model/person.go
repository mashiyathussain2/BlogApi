package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Person is the data structure that we will save and receive.
type Person struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	//Author_Id primitive.ObjectID     `json:"author_id,omitempty" bson:"author_id,omitempty"`
	FirstName    string                 `json:"first_name,omitempty" bson:"first_name,omitempty"`
	Comment_Info []Comment              `json:"comment_info" bson:"comment_info"`
	Blog_Info    []Blog                 `json:"blog_info" bson:"blog_info"`
	LastName     string                 `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Username     string                 `json:"username,omitempty" bson:"username,omitempty"`
	Email        string                 `json:"email,omitempty" bson:"email,omitempty"`
	Password     string                 `json:"password,omitempty" bson:"password,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"` // data is a optional fields that can hold anything in key:value format.
}
