package app

import (
	"encoding/json"
	"fmt"
	"time"

	"os"

	"log"
	"net/http"

	"blog/app/handler"
	"blog/app/schema"

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
	blogpage.Time = time.Now().UTC()
	result, err := db.Collection("blogpage").InsertOne(context.Background(), blogpage)
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
				"user_id":                         0,
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
				"post_id": 0,
			},
		},
	}
	groupStage := bson.D{
		{
			Key: "$group",
			Value: bson.M{
				"_id": "$_id",
				"title": bson.M{
					"$first": "$title",
				},
				"description": bson.M{
					"$first": "$description",
				},
				"blog_img": bson.M{
					"$first": "$blog_img",
				},
				"tag": bson.M{
					"$first": "$tag",
				},
				"category": bson.M{
					"$first": "$category",
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
	addfieldStage := bson.D{
		{
			Key: "$addFields",
			Value: bson.M{
				"time": bson.M{
					"$substr": bson.A{"$created_at", 0, 10},
				},
			},
		},
	}
	projectStage4 := bson.D{
		{
			Key: "$project",
			Value: bson.M{
				"created_at": 0,
			},
		},
	}
	// addfieldStage2 := bson.D{
	// 	{
	// 		Key: "$addFields",
	// 		Value: bson.M{
	// 			"month" : bson.M{
	// 				"$month": "$created_at",
	// 			},
	// 			"day" : bson.M{
	// 				"$dayOfMonth": "$created_at",
	// 			},
	// 			"year" : bson.M{
	// 				"$year": "$created_at",
	// 			},
	// 		},
	// 	},
	// }
	// addfieldStage3 := bson.D{
	// 	{
	// 		Key: "$addFields",
	// 		Value: bson.M{
	// 			"days" : bson.M{
	// 				"$toString": "$day",
	// 			},
	// 			"years" : bson.M{
	// 				"$toString": "$year",
	// 			},
	// 		},
	// 	},
	// }
	// projectStage5 := bson.D{
	// 	{
	// 		Key: "$project",
	// 		Value: bson.M{
	// 			"day": 0,
	// 			"year": 0,
	// 		},
	// 	},
	// }
	// addfieldStage4 := bson.D{
	// 	{
	// 		Key: "$addFields",
	// 		Value: bson.M{
	// 			"Month" : bson.M{
	// 				"$arrayElemAt" : [
	// 					[
	// 					  "",
	// 					  "Jan",
	// 					  "Feb",
	// 					  "Mar",
	// 					  "April",
	// 					  "May",
	// 					  "Jun",
	// 					  "Jul",
	// 					  "Aug",
	// 					  "Sep",
	// 					  "Oct",
	// 					  "Nov",
	// 					  "Dec",
	// 					],
	// 					"$month",
	// 				],

	// 			},
	// 		},
	// 	},
	// }

	pipeline := mongo.Pipeline{lookupStage, lookupStage2, projectStage, unwindStage, lookupStagesPeople, unwindStage2, projectStage2, unwindStage3, lookupStageLikes, projectStage3, groupStage, addfieldStage, projectStage4}

	// pageString := req.FormValue("page")
	// page, err := strconv.ParseInt(pageString, 10, 64)
	// if err != nil {
	// 	page = 0
	// }
	// page = page * limit
	// findOptions := options.FindOptions{
	// 	Skip:  &page,
	// 	Limit: &limit,
	// 	Sort: bson.M{
	// 		"_id": -1, // -1 for descending and 1 for ascending
	// 	},
	// }

	// // query for the aggregation
	showLoadedCursor, err := db.Collection("blogpage").Aggregate(context.TODO(), pipeline)
	if err != nil {
		fmt.Println("Hello", err)

	}
	var showsLoaded = []bson.M{}

	if err = showLoadedCursor.All(context.TODO(), &showsLoaded); err != nil {
		fmt.Println("Hellooo")

	}

	handler.ResponseWriter(res, http.StatusOK, "hello", showsLoaded)

}

// func SearchBlog(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
// 	var blogList []schema.Blog

// 	// query for find the user in the database
// 	curser, err := db.Collection("blogpage").Find(nil, bson.M{})
// 	if err != nil {
// 		log.Printf("Error while quering collection: %v\n", err)
// 		handler.ResponseWriter(res, http.StatusInternalServerError, "Error happend while reading data", nil)
// 		return
// 	}
// 	err = curser.All(context.Background(), &blogList)
// 	if err != nil {
// 		log.Fatalf("Error in curser: %v", err)
// 		handler.ResponseWriter(res, http.StatusInternalServerError, "Error happend while reading data", nil)
// 		return
// 	}
// 	handler.ResponseWriter(res, http.StatusOK, "", blogList)
// }

func GetBlog(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var params = mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "id that you sent is wrong!!!", nil)
		return
	}
	var blog bson.M

	err = db.Collection("blogpage").FindOne(context.Background(), bson.M{"_id": id}).Decode(&blog)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			handler.ResponseWriter(res, http.StatusNotFound, "blogpage not found", nil)
		default:
			log.Printf("Error while decode to go struct:%v\n", err)
			handler.ResponseWriter(res, http.StatusInternalServerError, "there is an error on server!!!", err.Error())
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
