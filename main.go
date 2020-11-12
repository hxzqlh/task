package main

import (
	"context"
	"task/common"
	"task/handler"
	task "task/proto/task"
	"task/repository"
	"task/subscriber"
	"time"

	"github.com/pkg/errors"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)

const MONGO_URI = "mongodb://admin:123456@127.0.0.1:27017"

func main() {
	conn, err := common.ConnectMongo(MONGO_URI, time.Second)
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
		micro.Name(common.TaskServiceName),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	taskHandler := &handler.TaskHandler{
		TaskRepository: &repository.TaskRepositoryImpl{
			Conn: conn,
		},
		TaskFinishedPubEvent: micro.NewEvent(common.TaskTopicName, service.Client()),
	}
	if err := task.RegisterTaskServiceHandler(service.Server(), taskHandler); err != nil {
		log.Fatal(errors.WithMessage(err, "register server"))
	}

	// Register Struct as Subscriber
	micro.RegisterSubscriber(common.TaskTopicName, service.Server(), new(subscriber.Task))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run server"))
	}
}
