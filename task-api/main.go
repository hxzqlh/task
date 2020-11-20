package main

import (
	"context"
	"task/common"
	pb "task/proto/task"

	"task/task-api/handler"
	"task/task-api/wrapper/breaker/hystrix"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	ot "github.com/opentracing/opentracing-go"
)

func main() {
	jaegerTracer, closer, err := common.NewJaegerTracer(common.TaskApiName, common.JaegerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()
	ot.SetGlobalTracer(jaegerTracer)

	etcdRegister := etcd.NewRegistry(registry.Addrs(common.EtcdAddr))

	app := micro.NewService(
		micro.Name(common.TaskClientName),
		micro.Registry(etcdRegister),
		micro.WrapClient(hystrix.NewClientWrapper(), opentracing.NewClientWrapper(jaegerTracer)),
	)

	cli := app.Client()
	cli.Init(
		client.Retries(3),
		client.Retry(func(ctx context.Context, req client.Request, retryCount int, err error) (bool, error) {
			log.Errorf("api retry call: %s.%s-%v", req.Service(), req.Method(), retryCount)
			return true, nil
		}),
	)

	taskService := pb.NewTaskService(common.TaskServiceName, cli)

	webHandler := gin.Default()
	// 这个服务才是真正运行的服务
	service := web.NewService(
		web.Name(common.TaskApiName),
		web.Address(common.TaskApiAddr),
		web.Handler(webHandler),
		web.Registry(etcdRegister),
	)
	// 配置web路由
	handler.Router(webHandler, taskService)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
