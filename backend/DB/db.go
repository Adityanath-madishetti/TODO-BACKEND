package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

//create a struct for the the Task and user here

var UserCollection *mongo.Collection
var Taskcollection *mongo.Collection


func MakeDbConnection()  *mongo.Client{


	conenctionString :="mongodb+srv://cs23btech11032:Aditya@5002@cluster0.hxuizul.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(conenctionString).SetServerAPIOptions(serverAPI)

	client , err := mongo.Connect(opts)

	if(err!=nil){
		log.Fatal("DB connection error:", err)
	}



	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second) // context is for tiem out while connecting db
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
        log.Fatal("❌ Ping failed:", err)
    }

	db := client.Database("TodoDatabase")
	UserCollection=db.Collection("Users")
	Taskcollection = db.Collection("TodoTasks")

	return client

}

