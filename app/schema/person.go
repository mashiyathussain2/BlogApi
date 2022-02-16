package schema

import (
	//"blog/app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Person is the data structure that we will save and receive.
type Person struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	//Author_Id primitive.ObjectID `json:"author_id,omitempty" bson:"author_id,omitempty"`
	FirstName string `json:"first_name,omitempty" bson:"first_name,omitempty" validate:"required"`
	//Person_Id primitive.ObjectID `json: "person_id,omitempty" bson:"person_id,omitempty"`
	LastName string `json:"last_name,omitempty" bson:"last_name,omitempty" validate:"required"`
	Email    string `json:"email,omitempty" bson:"email,omitempty" validate: "required,email"`
	//Blog_ID
	//Comment_Info map[string]Comment `json:"comment_info" bson:"comment_info"`
	//Comment_Info model.Comment      `json:"comment_info,omitempty" bson:"email,omitempty"`
	//model.Comment
	//FirstName string                 `json:"first_name,omitempty" bson:"first_name,omitempty"`
	//LastName  string                 `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	//Email     string                 `json:"email,omitempty" bson:"email,omitempty"`
	Data map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"` // data is a optional fields that can hold anything in key:value format.
}
