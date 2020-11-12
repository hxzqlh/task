package subscriber

import (
	"context"
	log "github.com/micro/go-micro/v2/logger"

	task "task/proto/task"
)

type Task struct{}

func (e *Task) Handle(ctx context.Context, msg *task.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *task.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}
