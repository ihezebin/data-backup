package service

import (
	"context"
	"data-backup/application/dto"
	"data-backup/component/task"

	"github.com/ihezebin/oneness/httpserver"
	"github.com/ihezebin/oneness/logger"
)

type RestoreApplicationService struct {
	logger *logger.Entry
}

func NewRestoreApplicationService(l *logger.Entry) *RestoreApplicationService {
	return &RestoreApplicationService{
		logger: l.WithField("application", "restore"),
	}
}

func (svc *RestoreApplicationService) RestoreByTask(ctx context.Context, req *dto.RestoreByTaskReq) (*dto.RestoreByTaskResp, error) {
	var _task *task.Task
	for _, t := range task.GetTasks() {
		if t.Id == req.TaskId {
			_task = t
			break
		}
	}

	if _task == nil {
		svc.logger.Errorf(ctx, "task not found, task_id: %s", req.TaskId)
		return nil, httpserver.ErrorWithBadRequest()
	}

	err := _task.Source.Restore(ctx, _task.Target)
	if err != nil {
		svc.logger.Errorf(ctx, "restore error, task: %+v, err: %v", _task, err)
		return nil, httpserver.NewError(httpserver.CodeInternalServerError, err.Error())
	}

	return &dto.RestoreByTaskResp{
		Task: _task,
	}, nil
}
