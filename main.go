package main

import (
	"context"
	"task/handler"
	"task/repository"
	"task/subscriber"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	task "task/proto/task"
)

const MONGO_URI = "mongodb://admin:123456@127.0.0.1:27017"

func main() {
	conn, err := connectMongo(MONGO_URI, time.Second)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = conn.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.task"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	taskHandler := &handler.TaskHandler{
		TaskRepository: &repository.TaskRepositoryImpl{
			Conn: conn,
		},
	}
	if err := task.RegisterTaskServiceHandler(service.Server(), taskHandler); err != nil {
		log.Fatal(errors.WithMessage(err, "register server"))
	}

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.task", service.Server(), new(subscriber.Task))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run server"))
	}
}

func connectMongo(uri string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
