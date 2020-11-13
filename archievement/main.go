package main

import (
	"context"
	"log"
	"task/archievement/repository"
	"task/archievement/subscriber"
	"task/common"
	"time"

	ot "github.com/opentracing/opentracing-go"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/broker/nats"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/pkg/errors"
)

// task-srv服务
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

	jaegerTracer, closer, err := common.NewJaegerTracer(common.ArchiementServiceName, common.JaegerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()
	ot.SetGlobalTracer(jaegerTracer)

	// New Service
	service := micro.NewService(
		micro.Name(common.ArchiementServiceName),
		micro.Version("latest"),
		micro.Registry(etcd.NewRegistry(registry.Addrs(common.EtcdAddr))),
		micro.Broker(nats.NewBroker(broker.Addrs(common.NatsAddr))),
		micro.WrapSubscriber(opentracing.NewSubscriberWrapper(jaegerTracer)),
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
