package handler

import (
	"context"
	"task/repository"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/pkg/errors"

	pb "task/proto/task"
)

type TaskHandler struct {
	TaskRepository       repository.TaskRepository
	TaskFinishedPubEvent micro.Event
}

func (t *TaskHandler) Create(ctx context.Context, req *pb.Task, resp *pb.EditResponse) error {
	log.Info("Received TaskHandler.Create request")
	if req.Body == "" || req.StartTime <= 0 || req.EndTime <= 0 || req.UserId == "" {
		return errors.New("bad param")
	}
	if err := t.TaskRepository.InsertOne(ctx, req); err != nil {
		return err
	}
	resp.Msg = "success"
	return nil
}

func (t *TaskHandler) Delete(ctx context.Context, req *pb.Task, resp *pb.EditResponse) error {
	log.Infof("Received TaskHandler.Delete request: %v", req.Id)
	if req.Id == "" {
		return errors.New("bad param")
	}
	if err := t.TaskRepository.Delete(ctx, req.Id); err != nil {
		return err
	}
	resp.Msg = req.Id
	return nil
}

func (t *TaskHandler) Modify(ctx context.Context, req *pb.Task, resp *pb.EditResponse) error {
	log.Infof("Received TaskHandler.Modify request: %v", req.Id)
	if req.Id == "" || req.Body == "" || req.StartTime <= 0 || req.EndTime <= 0 {
		return errors.New("bad param")
	}
	if err := t.TaskRepository.Modify(ctx, req); err != nil {
		return err
	}
	resp.Msg = "success"
	return nil
}

func (t *TaskHandler) Finished(ctx context.Context, req *pb.Task, resp *pb.EditResponse) error {
	log.Infof("Received TaskHandler.Finished request: %v", req.Id)
	if req.Id == "" || req.IsFinished != repository.UnFinished && req.IsFinished != repository.Finished {
		return errors.New("bad param")
	}
	if err := t.TaskRepository.Finished(ctx, req); err != nil {
		return err
	}
	resp.Msg = "success"

	// 发送task完成消息
	// 由于以下都是主业务之外的增强功能，出现异常只记录日志，不影响主业务返回
	if task, err := t.TaskRepository.FindById(ctx, req.Id); err != nil {
		log.Errorf("can't find finished task: %s", err.Error())
	} else {
		if err = t.TaskFinishedPubEvent.Publish(ctx, task); err != nil {
			log.Errorf("can't send task finished message: %s", err.Error())
		}
	}
	return nil
}

func (t *TaskHandler) Search(ctx context.Context, req *pb.SearchRequest, resp *pb.SearchResponse) error {
	log.Info("Received TaskHandler.Search request")
	count, err := t.TaskRepository.Count(ctx, req.Keyword)
	if err != nil {
		return errors.WithMessage(err, "count row number")
	}
	if req.PageCode <= 0 {
		req.PageCode = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.SortBy == "" {
		req.SortBy = "createTime"
	}
	if req.Order == 0 {
		req.Order = -1
	}
	if req.PageSize*(req.PageCode-1) > count {
		return errors.New("There's not that much data")
	}
	rows, err := t.TaskRepository.Search(ctx, req)
	if err != nil {
		return errors.WithMessage(err, "search data")
	}
	*resp = pb.SearchResponse{
		PageCode: req.PageCode,
		PageSize: req.PageSize,
		SortBy:   req.SortBy,
		Order:    req.Order,
		Rows:     rows,
	}
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *TaskHandler) Stream(ctx context.Context, req *pb.StreamingRequest, stream pb.TaskService_StreamStream) error {
	log.Infof("Received Task.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&pb.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *TaskHandler) PingPong(ctx context.Context, stream pb.TaskService_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&pb.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
