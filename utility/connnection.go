package utility

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func MongoConnection() {

	ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	Client, _ = mongo.Connect(ctx, clientOptions)

}
func DB() (*mongo.Collection, *mongo.Database) {
	MongoConnection()
	db := Client.Database("univercity")
	collection := Client.Database("univercity").Collection("Courseinfo")
	return collection, db
}
func DB1() *mongo.Collection {
	MongoConnection()
	collection := Client.Database("univercity").Collection("Studentinfo")
	return collection
}

// func DB2() *mongo.Collection {
// 	MongoConnection()
// 	collection := Client.Database("univercity").Collection("coursename")
// 	return collection
// }
