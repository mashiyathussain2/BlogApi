package schema

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Person is the data structure that we will save and receive.
type Person struct {
	ID        primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string                 `json:"first_name,omitempty" bson:"first_name,omitempty" validate:"required"`
	LastName  string                 `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     string                 `json:"email,omitempty" bson:"email,omitempty" validate: "required,email"`
	Password  string                 `json:"password,omitempty" bson:"password,omitmepty" validate:"required"`
	Username  string                 `json:"username,omitempty" bson:"username,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty" bson:"data,omitempty"` // data is a optional fields that can hold anything in key:value format.
}
