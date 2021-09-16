package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Collection {
	// set client options
	clietOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// connect to mongoDB
	client, err := mongo.Connect(context.TODO(), clietOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection := client.Database("users_db").Collection("users_info")

	return collection
}
