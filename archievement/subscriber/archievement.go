package subscriber

import (
	"context"
	"log"
	"strings"
	"task/archievement/repository"
	pb "task/proto/task"
	"time"

	"github.com/pkg/errors"
)

// 定时实现类
type AchievementSub struct {
	Repo repository.AchievementRepo
}

// 只处理任务完成这一个事件
func (sub *AchievementSub) Finished(ctx context.Context, task *pb.Task) error {
	log.Printf("Handler Received message: %v\n", task)
	if task.UserId == "" || strings.TrimSpace(task.UserId) == "" {
		return errors.New("userId is blank")
	}
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

func (sub *AchievementSub) Finished2(ctx context.Context, task *pb.Task) error {
	log.Println("Finished2")
	return nil
}

func (sub *AchievementSub) Finished3(ctx context.Context, task *pb.Task) error {
	log.Println("Finished3")
	return nil
}
