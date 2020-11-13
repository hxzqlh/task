package main

import (
	"log"
	"task/common"
	pb "task/proto/task"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
)

// task-srv服务的restful api映射
func main() {
	g := gin.Default()
	service := web.NewService(
		web.Name("go.micro.api.task"),
		web.Address(":8888"),
		web.Registry(etcd.NewRegistry(registry.Addrs(common.EtcdAddr))),
		web.Handler(g),
	)
	cli := pb.NewTaskService("go.micro.service.task", client.DefaultClient)

	v1 := g.Group("/task")
	{
		v1.GET("/search", func(c *gin.Context) {
			req := new(pb.SearchRequest)
			if err := c.BindQuery(req); err != nil {
				c.JSON(200, gin.H{
					"code": "500",
					"msg":  "bad param",
				})
				return
			}
			if resp, err := cli.Search(c, req); err != nil {
				c.JSON(200, gin.H{
					"code": "500",
					"msg":  err.Error(),
				})
			} else {
				c.JSON(200, gin.H{
					"code": "200",
					"data": resp,
				})
			}
		})
	}
	service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
