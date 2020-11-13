package main

import (
	"context"
	"log"
	"task/archievement/repository"
	"task/archievement/subscriber"
	"task/common"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/broker/nats"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/pkg/errors"
)

const MONGO_URI = "mongodb://admin:123456@127.0.0.1:27017"

// task-srv服务
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
		micro.Name("go.micro.service.achievement"),
		micro.Version("latest"),
		micro.Registry(etcd.NewRegistry(registry.Addrs(common.EtcdAddr))),
		micro.Broker(nats.NewBroker(broker.Addrs(common.NatsAddr))),
	)

	// Initialise service
	service.Init()

	// Register Handler
	handler := &subscriber.AchievementSub{
		Repo: &repository.AchievementRepoImpl{
			Conn: conn,
		},
	}
	// 这里的topic注意与task注册的要一致
	if err := micro.RegisterSubscriber(common.TaskTopicName, service.Server(), handler); err != nil {
		log.Fatal(errors.WithMessage(err, "subscribe"))
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(errors.WithMessage(err, "run server"))
	}
}
