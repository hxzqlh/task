package common

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MongodbUri = "mongodb://admin:123456@127.0.0.1:27017"
	EtcdAddr   = "127.0.0.1:2379"
)

const (
	TaskServiceName = "go.micro.service.task"
	TaskClientName  = "go.micro.client.task"
	TaskTopicName   = "go.micro.topic.task"
)

func ConnectMongo(uri string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
