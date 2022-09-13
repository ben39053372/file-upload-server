package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://ilovehk:password@mongo:27017/?maxPoolSize=20"

var database = "ilovehk-assets"

var collectionName = "assets"

var collection *mongo.Collection

func init() {

	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	// client = c

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
