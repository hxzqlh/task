package handler

import (
	"fmt"
	pb "task/proto/task"

	"github.com/gin-gonic/gin"
	ot "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
)

var service pb.TaskService

func Router(g *gin.Engine, taskService pb.TaskService) {
	service = taskService

	v1 := g.Group("/task")
	{
		v1.GET("/search", Search)
		v1.POST("/finished", Finished)
	}
}

func Search(c *gin.Context) {
	req := new(pb.SearchRequest)
	if err := c.BindQuery(req); err != nil {
		c.JSON(200, gin.H{
			"code": "500",
			"msg":  "bad param",
		})
		return
	}

	ctx, span, err := ot.StartSpanFromContext(c, opentracing.GlobalTracer(), c.Request.URL.Path)
	if err != nil {
		fmt.Println("start span err", err)
	}
	defer span.Finish()

	span.SetTag("req", req)

	if resp, err := service.Search(ctx, req); err != nil {
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
}

func Finished(c *gin.Context) {
	req := new(pb.Task)

	if err := c.BindJSON(req); err != nil {
		c.JSON(200, gin.H{
			"code": "500",
			"msg":  "bad param",
		})
		return
	}

	ctx, span, err := ot.StartSpanFromContext(c, opentracing.GlobalTracer(), c.Request.URL.Path)
	if err != nil {
		fmt.Println("start span err", err)
	}
	defer span.Finish()

	span.SetTag("req", req)
	span.LogFields(
		otlog.String("taskId", req.Id),
	)

	if resp, err := service.Finished(ctx, req); err != nil {
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
}
