package app

import (
	//"encoding/json"
	"log"
	"net/http"
	"strconv"

	"blog/app/handler"
	"blog/app/model"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// results count per page
//var limit int64 = 10

// CreatePerson will handle the create person post request
func CreateBlog(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	blog := new(model.Blog)
	blog.Title = req.FormValue("title")
	blog.Description = req.FormValue("description")
	blog.Author = req.FormValue("author")
	_, err := db.Collection("blogpage").InsertOne(context.TODO(), blog)
	//err := json.NewDecoder(req.Body).Decode(blogpage)

	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "body json request have issues!!!", nil)
		return
	}
	result, err := db.Collection("blogpage").InsertOne(nil, blog)
	//if err != nil {
	//	switch err.(type) {
	//	case mongo.WriteException:
	//		handler.ResponseWriter(res, http.StatusNotAcceptable, "username or email already exists in database.", nil)
	//	default:
	//		handler.ResponseWriter(res, http.StatusInternalServerError, "Error while inserting data.", nil)
	//	}
	//	return
	//}
	blog.ID = result.InsertedID.(primitive.ObjectID)
	handler.ResponseWriter(res, http.StatusCreated, "", blog)
}

// GetPersons will handle people list get request
func GetBlogs(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var blogpageList []model.Blog
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
	curser, err := db.Collection("blogpage").Find(nil, bson.M{}, &findOptions)
	if err != nil {
		log.Printf("Error while quering collection: %v\n", err)
		handler.ResponseWriter(res, http.StatusInternalServerError, "Error happend while reading data", nil)
		return
	}
	err = curser.All(context.Background(), &blogpageList)
	if err != nil {
		log.Fatalf("Error in curser: %v", err)
		handler.ResponseWriter(res, http.StatusInternalServerError, "Error happend while reading data", nil)
		return
	}
	handler.ResponseWriter(res, http.StatusOK, "", blogpageList)
}

// GetPerson will give us person with special id
func GetBlog(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var params = mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "id that you sent is wrong!!!", nil)
		return
	}
	var blogpage model.Blog
	err = db.Collection("blogpage").FindOne(nil, model.Blog{ID: id}).Decode(&blogpage)
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
	handler.ResponseWriter(res, http.StatusOK, "", blogpage)
}
