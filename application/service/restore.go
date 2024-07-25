package service

import (
	"context"
	"data-backup/application/dto"
	"data-backup/component/source"
	"data-backup/component/task"
	"path"

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
	var t *task.Task
	for _, _task := range task.GetTasks() {
		if _task.Id == req.TaskId {
			t = _task
			break
		}
	}

	if t == nil {
		svc.logger.Errorf(ctx, "task not found, task_id: %s", req.TaskId)
		return nil, httpserver.ErrorWithBadRequest()
	}

	_source := t.GetSource()
	if _source == nil {
		svc.logger.Errorf(ctx, "source not found, task: %+v", t)
		return nil, httpserver.ErrorWithInternalServer()
	}

	_target := t.GetTarget()
	if _target == nil {
		svc.logger.Errorf(ctx, "target not found, task: %+v", t)
		return nil, httpserver.ErrorWithInternalServer()
	}

	for _, sourceKey := range _source.Keys() {
		data, err := _target.Restore(ctx, sourceKey)
		if err != nil {
			svc.logger.Errorf(ctx, "restore error, source_key: %s, err: %v", sourceKey, err)
			return nil, httpserver.ErrorWithInternalServer()
		}
		if err = _source.Import(ctx, source.Data{
			Key:     path.Base(sourceKey),
			Content: data,
		}); err != nil {
			svc.logger.Errorf(ctx, "restore error, source_key: %s, err: %v", sourceKey, err)
			return nil, httpserver.ErrorWithInternalServer()
		}
	}

	return &dto.RestoreByTaskResp{
		Task: t,
	}, nil
}
