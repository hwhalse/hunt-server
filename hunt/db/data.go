// Package mongo extends the common mongo client, creating a connection for each service
package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

var HuntArtakClient = NewMongoClient("hunt")

func NewMongoClient(database string) *Client {
	client, err := NewClient(os.Getenv("MONGO_URI"), database)
	if err != nil {
		panic(err)
	}
	return &client
}

func CreateFilter(key string, val any) bson.M {
	return bson.M{key: val}
}
