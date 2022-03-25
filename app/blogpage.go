package app

import (
	"encoding/json"
	"fmt"
	"time"

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
	const (
		layoutISO = "2006-01-02"
		layoutUS  = "January 2, 2006"
	)
	//date := time.Now()
	t := time.Now().Format(layoutUS)
	result, err := db.Collection("blogpage").InsertOne(context.Background(), bson.M{"blogs": blogpage, "time": t})
	if err != nil {
		switch err.(type) {
		case mongo.WriteException:
			handler.ResponseWriter(res, http.StatusNotAcceptable, "username or email already exists in database.", nil)
		default:
			handler.ResponseWriter(res, http.StatusInternalServerError, "Error while inserting data.", nil)
		}
		return
	}
	// time.Parse
	tt := time.Now()
	//fmt.Println(t.Format("2006-01-02-15-04-05"))
	blogpage.ID = result.InsertedID.(primitive.ObjectID)
	blogpage.Time = tt.Format(layoutUS)
	handler.ResponseWriter(res, http.StatusCreated, "", blogpage)
}

// GetBlogs is for getting all blogs
func GetBlogs(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
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
	// unwindStage := bson.D{
	// 	{
	// 		Key: "$unwind",
	// 		Value: bson.M{
	// 			"path": "$comment",
	// 		},
	// 	},
	// }
	lookupStage2 := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         "people",
				"localField":   "blogs.user_id",
				"foreignField": "_id",
				"as":           "person_info",
			},
		},
	}
	projectStage := bson.D{
		{
			Key: "$project",
			Value: bson.M{
				"comment.post_id": 0,
			},
		},
	}
	unwindStage := bson.D{
		{
			Key: "$unwind",
			Value: bson.M{
				"path":                       "$comment",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	// unwindStage := bson.D{
	// 	{
	// 		Key: "$unwind",
	// 		Value: bson.M{
	// 			"path": "$comment",
	// 		},
	// 	},
	// }

	lookupStagesPeople := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         "people",
				"localField":   "comment.user_id",
				"foreignField": "_id",
				"as":           "comment.comment_author",
			},
		},
	}
	unwindStage2 := bson.D{
		{
			Key: "$unwind",
			Value: bson.M{
				"path":                       "$person_info",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}
	projectStage2 := bson.D{
		{
			Key: "$project",
			Value: bson.M{
				"comment.user_id":                 0,
				"comment.post_id":                 0,
				"comment.comment_author.password": 0,
				"comment.comment_author.email":    0,
				"person_info.password":            0,
				"person_info.email":               0,
				"blogs.user_id":                   0,
			},
		},
	}

	unwindStage3 := bson.D{
		{
			Key: "$unwind",
			Value: bson.M{
				"path":                       "$comment.comment_author",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}
	lookupStageLikes := bson.D{
		{
			Key: "$lookup",
			Value: bson.M{
				"from":         "like",
				"localField":   "_id",
				"foreignField": "post_id",
				"as":           "blog_likes",
			},
		},
	}
	projectStage3 := bson.D{
		{
			Key: "$project",
			Value: bson.M{
				"blog_likes.post_id": 0,
			},
		},
	}
	groupStage := bson.D{
		{
			Key: "$group",
			Value: bson.M{
				"_id": "$_id",
				"title": bson.M{
					"$first": "$blogs.title",
				},
				"description": bson.M{
					"$first": "$blogs.description",
				},
				"blog_img": bson.M{
					"$first": "$blogs.blog_img",
				},
				"comment": bson.M{
					"$push": "$comment",
				},
				"author_info": bson.M{
					"$first": "$person_info",
				},
				"likes": bson.M{
					"$first": "$blog_likes",
				},
				"created_at": bson.M{
					"$first": "$time",
				},
			},
		},
	}
	// likeLookup := bson.D{
	// 	{
	// 		Key: "$lookup",
	// 		Value: bson.M{
	// 			"from":         "like",
	// 			"localField":   "_id",
	// 			"foreignField": "post_id",
	// 			"as":           "blog_likes",
	// 		},
	// 	},
	// }

	pipeline := mongo.Pipeline{lookupStage /*unwindStage,*/, lookupStage2 /* unwindStage,, lookupStagesPeople*/ /*unwindStageCommentAuthor*/, projectStage, unwindStage, lookupStagesPeople, unwindStage2, projectStage2, unwindStage3, lookupStageLikes, projectStage3, groupStage}

	// // query for the aggregation
	// showLoadedCursor, err := db.Collection("blogpage").Aggregate(context.TODO(), pipeline)
	showLoadedCursor, err := db.Collection("blogpage").Aggregate(context.TODO(), pipeline)
	if err != nil {
		fmt.Println("Hello", err)

	}
	var showsLoaded = []bson.M{}
	//showsLoaded := new(schema.Blog)

	if err = showLoadedCursor.All(context.TODO(), &showsLoaded); err != nil {
		fmt.Println("Hellooo")

	}
	// count, err := db.Collection("blogpage").CountDocuments(context.TODO(), bson.M{})
	// fmt.Println(count, err)
	//now := time.Now()
	//fmt.Println(showsLoaded)

	handler.ResponseWriter(res, http.StatusOK, "hello", showsLoaded)

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
