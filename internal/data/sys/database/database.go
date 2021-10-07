package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	URI string
}

func Open(cfg Config) (*mongo.Client, error) {
	clientOptions := options.Client().
		ApplyURI(cfg.URI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer client.Disconnect(ctx)
	log.Println("connected to Database")
	return client, nil
}
