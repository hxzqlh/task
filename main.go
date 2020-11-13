package main

import (
	"context"
	"task/common"
	"task/handler"
	task "task/proto/task"
	"task/repository"
	"task/subscriber"
	"time"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/pkg/errors"
)

func main() {
	conn, err := common.ConnectMongo(common.MongodbUri, time.Second)
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
		micro.Registry(etcd.NewRegistry(registry.Addrs(common.EtcdAddr))),
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
