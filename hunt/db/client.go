package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type Client struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewClient(uri, dbName string) (Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().
    ApplyURI(uri).
    SetMaxPoolSize(500).
    SetMinPoolSize(10)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return Client{}, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		// Panic needed because if DB fails, app should crash
		panic(err)
	}
	return Client{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

func (c *Client) Database() *mongo.Database {
	return c.db
}
