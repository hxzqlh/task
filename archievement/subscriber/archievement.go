package subscriber

import (
	"context"
	"fmt"
	"log"
	"strings"
	"task/archievement/repository"
	"task/common"
	pb "task/proto/task"
	"time"

	ot "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type AchievementSub struct {
	Repo repository.AchievementRepo
}

// 只处理任务完成这一个事件
func (sub *AchievementSub) Finished(c context.Context, task *pb.Task) error {
	log.Printf("Handler Received message: %v\n", task)
	if task.UserId == "" || strings.TrimSpace(task.UserId) == "" {
		return errors.New("userId is blank")
	}

	ctx, span, err := ot.StartSpanFromContext(c, opentracing.GlobalTracer(), common.ArchiementServiceName+".OnFinished")
	if err != nil {
		fmt.Println("start span err", err)
	}
	defer span.Finish()

	span.LogFields(
		otlog.String("taskId", task.Id),
	)

	entity, err := sub.Repo.FindByUserId(ctx, task.UserId)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	if entity == nil {
		entity = &repository.Achievement{
			UserId:        task.UserId,
			Total:         1,
			Finished1Time: now,
		}
		return sub.Repo.Insert(ctx, entity)
	}
	entity.Total++
	switch entity.Total {
	case 100:
		entity.Finished100Time = now
	case 1000:
		entity.Finished1000Time = now
	}
	return sub.Repo.Update(ctx, entity)

}
