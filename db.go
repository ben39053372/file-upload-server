package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://ilovehk:password@mongo:27017/?maxPoolSize=20"

var Client *mongo.Client

type FilesData struct {
	Id        primitive.ObjectID
	Size      int64     `bson:"size,omitempty"`
	CreatedAt time.Time `bson:"createdAt,omitempty"`
}

func init() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	Client = client

}

func Insert(data FilesData) *mongo.InsertOneResult {
	coll := Client.Database("db").Collection("assets")
	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		log.Panic(err)
	}
	return result
}

func Get(id primitive.ObjectID) bson.M {
	var result bson.M
	coll := Client.Database("db").Collection("assets")
	coll.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)

	return result
}
