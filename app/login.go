package app

import (
	"blog/app/handler"
	"blog/app/helpers"
	"blog/app/model"
	"blog/app/schema"
	"encoding/json"

	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/net/context"
)

// Login is the post request for login the person
func Login(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	user := schema.Login{}
	var foundUser model.Person
	// decode the user details
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&user)
	if err != nil {
		panic(err)
	}
	// query for finding the user in the database
	err = db.Collection("people").FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		handler.ResponseWriter(res, http.StatusInternalServerError, "login or passowrd is incorrect", err)
		return
	}
	// verify the user password
	passwordIsValid, msg := VerifyPassword(user.Password, foundUser.Password)
	if passwordIsValid != true {
		handler.ResponseWriter(res, http.StatusInternalServerError, msg, passwordIsValid)
		return
	}
	// verify user email
	if &user.Email == (*string)(nil) {
		handler.ResponseWriter(res, http.StatusInternalServerError, "user not found", nil)
		return
	}
	// generate the token on login.
	token, err := helpers.GenerateAllTokens(*&foundUser.Email, *&foundUser.Password)
	// unmarshal the token value.
	err = json.Unmarshal([]byte(token), err)

	err = db.Collection("people").FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		handler.ResponseWriter(res, http.StatusInternalServerError, "", nil)
		return
	}
	// save that token in cookie.
	http.SetCookie(res,
		&http.Cookie{
			Name:  "token",
			Value: token,
		})
	// return the token whenever success login.
	handler.ResponseWriter(res, http.StatusOK, "Success Login", token)
}
