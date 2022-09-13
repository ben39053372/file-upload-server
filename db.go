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

var database = "ilovehk-assets"

var collectionName = "assets"

var collection *mongo.Collection

type FilesData struct {
	Size      int64     `bson:"size,omitempty"`
	CreatedAt time.Time `bson:"createdAt,omitempty"`
	FilePath  string    `bson:"filePath,omitempty"`
	Url       string    `bson:"url,omitempty"`
}

func init() {

	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	// client = c

	collection = c.Database(database).Collection(collectionName)

}

func Insert(data FilesData) *mongo.InsertOneResult {
	result, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		log.Panic(err)
	}
	return result
}

func CreateEmptyDoc() (interface{}, error) {
	result, err := collection.InsertOne(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	newId := result.InsertedID
	return newId, nil
}

func UpdateDoc(id interface{}, data FilesData) (bson.M, error) {
	var result bson.M
	err := collection.FindOneAndUpdate(
		context.Background(),
		bson.D{{"_id", id}},
		bson.D{{"$set", data}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&result)
	return result, err
}

func Get(id primitive.ObjectID) bson.M {
	var result bson.M
	collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)

	return result
}
