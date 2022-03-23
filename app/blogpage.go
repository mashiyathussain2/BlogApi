package app

import (
	"encoding/json"
	"fmt"

	"os"

	"log"
	"net/http"

	"blog/app/handler"
	"blog/app/helpers"
	"blog/app/schema"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var jwtKey string = os.Getenv("SECRET_KEY")

// CreateBlog will handle the create blog post request
func CreateBlog(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	blogpage := new(schema.Blog)
	err := json.NewDecoder(req.Body).Decode(blogpage)
	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "body json request have issues!!!", nil)
		return
	}
	// for checking authorization
	cookie, err := req.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			handler.ResponseWriter(res, http.StatusUnauthorized, "Unauthorized", err.Error()) //res.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler.ResponseWriter(res, http.StatusUnauthorized, "Unauthorized", nil) //res.WriteHeader(http.StatusBadRequest)
		return
	}
	// tokenstr as the value of cookie
	tokenStr := cookie.Value

	claims := helpers.SignedDetails{}

	tkn, err := jwt.ParseWithClaims(tokenStr, &claims,
		func(tkn *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			handler.ResponseWriter(res, http.StatusUnauthorized, "Unauthorized", nil) //res.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler.ResponseWriter(res, http.StatusBadRequest, "Bad Request", err.Error()) //res.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		handler.ResponseWriter(res, http.StatusUnauthorized, "Unauthorized", nil) //res.WriteHeader(http.StatusUnauthorized)
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
	handler.ResponseWriter(res, http.StatusCreated, "", blogpage)
}

// GetBlogs is for getting all blogs
func GetBlogs(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	//var showsLoaded []schema.Blog
	// aggregation method starts from here.
	lookupStage := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         "comment",
				"localField":   "_id",
				"foreignField": "post_id",
				"as":           "comment",
			},
		},
	}

	lookupStage2 := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         "people",
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "person_info",
			},
		},
	}

	unwindStage := bson.D{
		{
			Key: "$unwind",
			Value: bson.M{
				"path": "$comment",
			},
		},
	}

	lookupStagesPeople := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         "people",
				"localField":   "comment.user_id",
				"foreignField": "_id",
				"as":           "comment.author_info",
			},
		},
	}

	unwindStageCommentAuthor := bson.D{
		{
			Key: "$unwind",
			Value: bson.M{
				"path": "$comment.author_info",
			},
		},
	}
	groupStage := bson.D{
		{
			Key: "$group",
			Value: bson.M{
				"_id": "$_id",
				"description": bson.M{
					"$first": "$description",
				},
				"comment": bson.M{
					"$push": "$comment",
				},
			},
		},
	}
	likeLookup := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         "like",
				"localField":   "_id",
				"foreignField": "post_id",
				"as":           "likes",
			},
		},
	}

	pipeline := mongo.Pipeline{lookupStage, lookupStage2, unwindStage, lookupStagesPeople, unwindStageCommentAuthor, groupStage, likeLookup}
	// query for the aggregation
	showLoadedCursor, err := db.Collection("blogpage").Aggregate(context.Background(), pipeline)
	if err != nil {
		fmt.Println("1", err)
		return
	}
	showsLoaded := new([]schema.Blog)

	if err = showLoadedCursor.All(context.TODO(), &showsLoaded); err != nil {
		fmt.Println("2", err)
	}

	count, err := db.Collection("blogpage").CountDocuments(context.TODO(), bson.M{})
	fmt.Println(count, err)

	fmt.Println(showsLoaded)
	handler.ResponseWriter(res, http.StatusOK, "", showsLoaded)

}

// GetBlog will give us blog with special id
func GetBlog(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var params = mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "id that you sent is wrong!!!", nil)
		return
	}
	var blog schema.Blog
	err = db.Collection("blogpage").FindOne(nil, schema.Blog{ID: id}).Decode(&blog)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			handler.ResponseWriter(res, http.StatusNotFound, "blog not found", err.Error())
		default:
			log.Printf("Error while decode to go struct:%v\n", err)
			handler.ResponseWriter(res, http.StatusInternalServerError, "there is an error on server!!!", nil)
		}
		return
	}
	handler.ResponseWriter(res, http.StatusOK, "", blog)
}

// UpdateBlog will handle the blog update endpoint
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
	result, err := db.Collection("blogpage").UpdateOne(context.Background(), schema.Blog{ID: oid}, update)
	if err != nil {
		log.Printf("Error while updateing document: %v", err)
		handler.ResponseWriter(res, http.StatusInternalServerError, "error in updating document!!!", nil)
		return
	}
	if result.MatchedCount == 1 {
		handler.ResponseWriter(res, http.StatusAccepted, "", &updateData)
	} else {
		handler.ResponseWriter(res, http.StatusNotFound, "blog not found", nil)
	}
}
