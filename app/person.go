package app

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"blog/app/handler"
	"blog/app/schema"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"

	//uniquevalidator "github.com/ramadani/go-unique-validator"

	//uniquevalidator "github.com/ramadani/go-unique-validator"
	//"github.com/thedevsaddam/govalidator"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// results count per page
var limit int64 = 10
var validate *validator.Validate
var uni *ut.UniversalTranslator

//var coll *mongo.Collection

// CreatePerson will handle the create person post request
func CreatePerson(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	en := en.New()
	uni = ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	person := new(schema.Person)
	err := json.NewDecoder(req.Body).Decode(person)
	err = validate.Struct(person)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			errs := err.(validator.ValidationErrors)
			fmt.Println(errs.Translate(trans))
			//fmt.Println("\n%+v\n", err)
			return
		}
	}

	//var coll *mongo.Collection
	//var id primitive.ObjectID

	// find the document for which the _id field matches id
	// specify the Sort option to sort the documents by age
	// the first document in the sorted order will be returned
	//opts := options.FindOne().SetSort(bson.D{{"_id", 1}})
	//var results bson.M
	//var person schema.Person
	err = db.Collection("people").FindOne(nil, schema.Person{Email: person.Email}).Decode(&person)
	if err != nil {
		result, err := db.Collection("people").InsertOne(context.Background(), person)
		if err != nil {
			switch err.(type) {
			case mongo.WriteException:
				handler.ResponseWriter(res, http.StatusNotAcceptable, "username or email already exists in database.", nil)
			default:
				handler.ResponseWriter(res, http.StatusInternalServerError, "Error while inserting data.", nil)
			}
			return
		}
		person.ID = result.InsertedID.(primitive.ObjectID)
		//person.Person_Id = result.InsertedID.(primitive.ObjectID)
		handler.ResponseWriter(res, http.StatusCreated, "", person)
	}
	handler.ResponseWriter(res, http.StatusNotAcceptable, "username or email already exists in database.", nil)

	//fmt.Printf("found document %v", person)

	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "body json request have issues!!!", err.Error())
		return
	}
	//var e mongo.WriteException
	//if err != nil {
	//	for _, we := range e.WriteErrors {
	//		if we.Code == 11000 {
	//			handler.ResponseWriter(res, http.StatusNotAcceptable, "username or email already exists in database.", err.Error())
	//		}
	//	}
	//}
	//uniqueRule := uniquevalidator.NewUniqueRule(db, "unique")
	//govalidator.AddCustomRule("unique", uniqueRule.Rule)
	//err := db.Collection.ensureIndex( "email", unique )
	//result, err := db.Collection("people").InsertOne(req.Context(), person)
	//err := db.Collection.createIndex( email: 1 , unique: true )
	//err = db.collection.createIndex( { "email": 1 }, { unique: true } )

}

// GetPersons will handle people list get request
func GetPersons(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var personList []schema.Person
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

	curser, err := db.Collection("people").Find(nil, bson.M{}, &findOptions)
	if err != nil {
		log.Printf("Error while quering collection: %v\n", err)
		handler.ResponseWriter(res, http.StatusInternalServerError, "Error happend while reading data", nil)
		return
	}
	err = curser.All(context.Background(), &personList)
	if err != nil {
		log.Fatalf("Error in curser: %v", err)
		handler.ResponseWriter(res, http.StatusInternalServerError, "Error happend while reading data", nil)
		return
	}
	handler.ResponseWriter(res, http.StatusOK, "", personList)
}

// GetPerson will give us person with special id
func GetPerson(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var params = mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		handler.ResponseWriter(res, http.StatusBadRequest, "id that you sent is wrong!!!", nil)
		return
	}
	var person schema.Person
	err = db.Collection("people").FindOne(nil, schema.Person{ID: id}).Decode(&person)
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
	handler.ResponseWriter(res, http.StatusOK, "", person)
}

// UpdatePerson will handle the person update endpoint
func UpdatePerson(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
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
	result, err := db.Collection("people").UpdateOne(context.Background(), schema.Person{ID: oid}, update)
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
