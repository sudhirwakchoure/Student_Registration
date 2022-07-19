package main

import (
	"STUDENT_REGISTRATION/controller"
	"STUDENT_REGISTRATION/utility"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	fmt.Println("Go  Gin framework")
	ctx, canc := context.WithTimeout(context.Background(), 10*time.Second)
	utility.Addindex()
	defer canc()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	utility.Client, _ = mongo.Connect(ctx, clientOptions)

	fmt.Println("connected")
	controller.HandleRequests()

}
