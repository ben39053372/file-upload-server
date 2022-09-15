package main

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri = ""

var database = "ben39053372/file-server"

var collectionName = "assets-data"

var collection *mongo.Collection

func init() {
	uri = os.Getenv("DB_URI")
	if uri == "" {
		panic("missing env: DB_URI")
	}
	print("uri:" + uri)

	if database := os.Getenv("DB_NAME"); database == "" {
		database = "file-server"
	}

	if collectionName := os.Getenv("DB_COLLECTION"); collectionName == "" {
		collectionName = "assets-data"
	}

	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	collection = c.Database(database).Collection(collectionName)

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
		bson.D{{Key: "_id", Value: id}},
		bson.D{{Key: "$set", Value: data}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&result)
	return result, err
}

func Get(id primitive.ObjectID) (FilesData, error) {
	var result FilesData
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)

	return result, err
}

func Del(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	result, err := collection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: id}})
	return result, err
}
