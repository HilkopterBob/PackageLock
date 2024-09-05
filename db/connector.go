package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client mongo.Client
	DB     *mongo.Database
)

func ConnectDb(connectionURI string) (mongo.Client, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))
	if err != nil {
		// TODO: log
		return mongo.Client{}, err
	}

	return *client, nil
}
