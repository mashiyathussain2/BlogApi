package app

import (
	"encoding/json"
	"log"
	"net/http"

	//"strconv"

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
	blogpage := new(model.Blog)
	err := json.NewDecoder(req.Body).Decode(blogpage)
	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "body json request have issues!!!", nil)
		return
	}
	result, err := db.Collection("blogpage").InsertOne(nil, blogpage)
	if err != nil {
		switch err.(type) {
		case mongo.WriteException:
			handler.ResponseWriter(res, http.StatusNotAcceptable, "username or email already exists in database.", nil)
		default:
			handler.ResponseWriter(res, http.StatusInternalServerError, "Error while inserting data.", nil)
		}
		return
	}
	blogpage.ID = result.InsertedID.(primitive.ObjectID)
	blogpage.Author_Id = result.InsertedID.(primitive.ObjectID)
	handler.ResponseWriter(res, http.StatusCreated, "", blogpage)
}

// GetBlogs is for getting all blogs
func GetBlogs(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var blogs []*model.Blog
	cur, err := db.Collection("blogpage").Find(context.TODO(), bson.M{}, options.Find())
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem model.Blog
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		blogs = append(blogs, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())

	// trying to lookup
	//lookupStage := bson.M{{"$lookup", bson.M{{"from", "podcasts"}, {"localField", "podcast"}, {"foreignField", "_id"}, {"as", "podcast"}}}}

	lookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         "person",
				"localField":   "_id",
				"foreignField": "author_id",
				"as":           "author_info",
			},
		},
	}
	showLoadedCursor, eerr := db.Collection("blogpage").Aggregate(context.TODO(), mongo.Pipeline{lookupStage})

	if eerr != nil {
		log.Fatal(eerr)
	}
	var showsLoaded []bson.M
	if eerr = showLoadedCursor.All(context.TODO(), &showsLoaded); err != nil {
		log.Fatal(eerr)
	}
	//_, eerr := db.Collection("people").Aggregate(context.TODO(),bson.M{
	//	"$lookup" : bson.M{
	//		"from" : "person",
	//		"localField" : "_id",
	//		"foreignField" : "author_id",
	//		"as" : "author_info",
	//	}});
	//if eerr != nil {
	//	log.Fatal(eerr)
	//}

	handler.ResponseWriter(res, http.StatusOK, "", blogs)
	//respondJSON(res, http.StatusOK, blogs)
}

// GetPerson will give us person with special id
func GetBlog(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var params = mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "id that you sent is wrong!!!", nil)
		return
	}
	var blog model.Blog
	err = db.Collection("people").FindOne(nil, model.Blog{ID: id}).Decode(&blog)
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
	handler.ResponseWriter(res, http.StatusOK, "", blog)
}

// UpdatePerson will handle the person update endpoint
func UpdateBlog(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var updateData map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&updateData)
	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "json body is incorrect", nil)
		return
	}
	// we dont handle the json decode return error because all our fields have the omitempty tag.
	var params = mux.Vars(req)
	oid, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "id that you sent is wrong!!!", nil)
		return
	}
	update := bson.M{
		"$set": updateData,
	}
	result, err := db.Collection("people").UpdateOne(context.Background(), model.Blog{ID: oid}, update)
	if err != nil {
		log.Printf("Error while updateing document: %v", err)
		handler.ResponseWriter(res, http.StatusInternalServerError, "error in updating document!!!", nil)
		return
	}
	if result.MatchedCount == 1 {
		handler.ResponseWriter(res, http.StatusAccepted, "", &updateData)
	} else {
		handler.ResponseWriter(res, http.StatusNotFound, "person not found", nil)
	}
}
