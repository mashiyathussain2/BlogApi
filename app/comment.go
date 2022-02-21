package app

import (
	"encoding/json"
	"log"
	"strconv"

	"net/http"

	"blog/app/handler"
	"blog/app/schema"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// CreateComment will handle the create comment post request
func CreateComment(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	comment := new(schema.Comment)
	err := json.NewDecoder(req.Body).Decode(comment)
	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "body json request have issues!!!", nil)
		return
	}
	// query for insert one comment in the database.
	result, err := db.Collection("comment").InsertOne(nil, comment)
	if err != nil {
		switch err.(type) {
		case mongo.WriteException:
			handler.ResponseWriter(res, http.StatusNotAcceptable, "username or email already exists in database.", nil)
		default:
			handler.ResponseWriter(res, http.StatusInternalServerError, "Error while inserting data.", nil)
		}
		return
	}
	comment.ID = result.InsertedID.(primitive.ObjectID)
	handler.ResponseWriter(res, http.StatusCreated, "", comment)
}

// GetComments will handle the get comment get request
func GetComments(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var commentList []schema.Comment
	pageString := req.FormValue("page")
	page, err := strconv.ParseInt(pageString, 10, 64)
	if err != nil {
		page = 0
	}
	page = page * limit
	findOptions := options.FindOptions{
		Skip:  &page,
		Limit: &limit,
		Sort: bson.M{
			"_id": -1, // -1 for descending and 1 for ascending
		},
	}
	// query for find the comment.
	curser, err := db.Collection("comment").Find(nil, bson.M{}, &findOptions)
	if err != nil {
		log.Printf("Error while quering collection: %v\n", err)
		handler.ResponseWriter(res, http.StatusInternalServerError, "Error happend while reading data", nil)
		return
	}
	// return all commentlist.
	err = curser.All(context.Background(), &commentList)
	if err != nil {
		log.Fatalf("Error in curser: %v", err)
		handler.ResponseWriter(res, http.StatusInternalServerError, "Error happend while reading data", nil)
		return
	}
	handler.ResponseWriter(res, http.StatusOK, "", commentList)
}

// GetComment will give us comment with special id
func GetComment(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var params = mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "id that you sent is wrong!!!", nil)
		return
	}
	var comment schema.Comment
	// query for find one comment in the database.
	err = db.Collection("comment").FindOne(nil, schema.Comment{ID: id}).Decode(&comment)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			handler.ResponseWriter(res, http.StatusNotFound, "person not found", nil)
		default:
			log.Printf("Error while decode to go struct:%v\n", err)
			handler.ResponseWriter(res, http.StatusInternalServerError, "there is an error on server!!!", nil)
		}
		return
	}
	handler.ResponseWriter(res, http.StatusOK, "", comment)
}
