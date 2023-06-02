package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collections DBCollections

type DBWorker struct {
	client *mongo.Client
}

type DBCollections struct {
	Motel mongo.Collection
}

func NewDBWorker(c *mongo.Client) *DBWorker {
	return &DBWorker{
		client: c,
	}
}

func (cfw *DBWorker) start() {
	Collections = DBCollections{
		Motel: *cfw.client.Database("test").Collection("motels"),
	}
}

func RunDB() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:mongolin@localhost:27017"))
	if err != nil {
		panic(err)
	}
	worker := NewDBWorker(client)
	worker.start()
}
