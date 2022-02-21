package schema

// Login is the data structure that we will save and receive.
type Login struct {
	Email    string `json:"email,omitempty" bson:"email,omitempty" validate: "required,email"`
	Password string `json:"password,omitempty" bson:"password,omitmepty" validate:"required"`
}
