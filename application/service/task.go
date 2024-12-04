package service

import (
	"context"
	"data-backup/application/dto"
	"data-backup/component/task"

	"github.com/ihezebin/oneness/httpserver"
	"github.com/ihezebin/oneness/logger"
)

type TaskApplicationService struct {
	logger *logger.Entry
}

func NewTaskApplicationService(l *logger.Entry) *TaskApplicationService {
	return &TaskApplicationService{
		logger: l.WithField("application", "task"),
	}
}

func (svc *TaskApplicationService) Restore(ctx context.Context, req *dto.TaskRestoreReq) (*dto.TaskRestoreResp, error) {
	taskIdM := make(map[string]struct{})
	for _, id := range req.TaskIds {
		taskIdM[id] = struct{}{}
	}

	if req.TaskId != "" {
		taskIdM[req.TaskId] = struct{}{}
	}

	tasks := make([]*task.Task, 0)
	for _, t := range task.GetTasks() {
		if _, ok := taskIdM[t.Id]; ok {
			tasks = append(tasks, t)
		}
	}

	for _, _task := range tasks {
		if _task == nil {
			svc.logger.Errorf(ctx, "task not found, task_id: %s", req.TaskId)
			return nil, httpserver.ErrorWithBadRequest()
		}

		err := _task.Source.Restore(ctx, _task.Target)
		if err != nil {
			svc.logger.WithError(err).Errorf(ctx, "task error, task: %+v", _task)
			return nil, httpserver.NewError(httpserver.CodeInternalServerError, err.Error())
		}
	}

	return &dto.TaskRestoreResp{
		Tasks: tasks,
	}, nil
}

func (svc *TaskApplicationService) Trigger(ctx context.Context, req *dto.TaskTriggerReq) (*dto.TaskTriggerResp, error) {
	for _, _task := range task.GetTasks() {
		if _task.Id == req.TaskId {
			err := _task.Source.Backup(ctx, _task.Target)
			if err != nil {
				svc.logger.WithError(err).Errorf(ctx, "task error, task: %+v", _task)
				return nil, httpserver.NewError(httpserver.CodeInternalServerError, err.Error())
			}

			return &dto.TaskTriggerResp{
				Task: _task,
			}, nil
		}

	}
	svc.logger.Errorf(ctx, "task not found, task_id: %s", req.TaskId)
	return nil, httpserver.ErrorWithBadRequest()
}
