package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(uri string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoDB{
		Client:   client,
		Database: client.Database("quckapp_media"),
	}, nil
}

func (m *MongoDB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	m.Client.Disconnect(ctx)
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

func (m *MongoDB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return m.Client.Ping(ctx, nil)
}
