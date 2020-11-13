package main

import (
	"log"
	"task/common"
	pb "task/proto/task"

	"task/task-api/handler"
	"task/task-api/wrapper/breaker/hystrix"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
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
	taskService := pb.NewTaskService(common.TaskServiceName, app.Client())

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
