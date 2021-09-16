package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sajid/data"
	"github.com/sajid/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collection = db.ConnectDB()

func GetAllUsers(response http.ResponseWriter, request *http.Request) {
	println("handle get all users request")

	// create user slice
	var users []data.User

	// passing empty filter to get all data
	// bson.M{} - unordered
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		http.Error(response, err.Error(), 500)
		return
	}

	for cur.Next(context.TODO()) {
		// create a user where a single document is decoded
		var user data.User

		err := cur.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		if user.Active {
			users = append(users, user)
		}
	}

	err = cur.Err()
	if err != nil {
		log.Fatal(err)
	}

	// encoding users to response
	encoder := json.NewEncoder(response)
	encoder.Encode(users)
}

func GetUser(response http.ResponseWriter, request *http.Request) {
	println("handle get a user request")

	var user data.User
	id, _ := getIdFromRequest(request)
	// create filter document
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		http.Error(response, err.Error(), http.StatusNotFound)
		return
	}

	if !user.Active {
		http.Error(response, "User does not exist", 404)
		return
	}

	encoder := json.NewEncoder(response)
	encoder.Encode(user)
}

func CreateUser(response http.ResponseWriter, request *http.Request) {
	println("handle create user")

	var user data.User
	json.NewDecoder(request.Body).Decode(&user)
	user.Validate()
	user.Created = time.Now()
	user.Updated = time.Now()
	user.Active = true
	result, err := collection.InsertOne(context.TODO(), user)

	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(response).Encode(result)
}

func UpdateUser(response http.ResponseWriter, request *http.Request) {
	println("handle update user")

	id, _ := getIdFromRequest(request)

	// read updated user from request body
	var user data.User
	json.NewDecoder(request.Body).Decode(&user)
	user.Validate()

	// prepare updated model
	update := bson.D{
		{"$set", bson.D{
			{"updated", time.Now()},
			{"firstName", user.FirstName},
			{"lastName", user.LastName},
			{"age", bson.D{
				{"value", user.Age.Value},
				{"interval", user.Age.Interval},
			}},
			{"mobile", user.Mobile},
		}},
	}

	// create filter document
	filter := bson.M{"_id": id}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&user)

	if err != nil {
		db.GetError(err, response)
		return
	}

	json.NewEncoder(response).Encode(user)
}

func DeleteUser(response http.ResponseWriter, request *http.Request) {
	println("handle delete request")

	id, _ := getIdFromRequest(request)
	// filter document
	filter := bson.M{"_id": id}

	update := bson.D{
		{"$set", bson.D{
			{"active", false},
		}},
	}

	// getting the user
	var user data.User
	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&user)

	if err != nil {
		http.Error(response, err.Error(), http.StatusNotFound)
		return
	}
}

func getIdFromRequest(request *http.Request) (primitive.ObjectID, error) {
	params := mux.Vars(request)
	return primitive.ObjectIDFromHex(params["id"])
}
