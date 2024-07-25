package cron

import (
	"context"
	"data-backup/component/email"
	"data-backup/component/task"
	"data-backup/config"
	"fmt"
	"time"

	mail "github.com/ihezebin/oneness/email"
	"github.com/ihezebin/oneness/logger"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

func Run(ctx context.Context) error {
	c := cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))

	tasks := task.GetTasks()
	for _, _task := range tasks {
		if _, err := c.AddFunc(_task.Cron, func() {
			var err error
			defer func() {
				if err != nil {
					logger.WithError(err).Errorf(ctx, "task error, task: %+v", _task)
					if nErr := notice(ctx, _task, err); nErr != nil {
						logger.WithError(nErr).Errorf(ctx, "task notice error, task: %+v", _task)
					}
				}
			}()

			logger.Infof(ctx, "task start, task: %+v", _task)

			source := _task.GetSource()
			if source == nil {
				err = errors.Errorf("task source is nil")
				return
			}

			target := _task.GetTarget()
			if target == nil {
				err = errors.Errorf("task target is nil")
				return
			}

			data, err := source.Export(ctx)
			if err != nil {
				err = errors.Wrapf(err, "source export error")
				return
			}

			for _, datum := range data {
				err = target.Store(ctx, datum.Key, datum.Content)
				if err != nil {
					err = errors.Wrapf(err, "target store error, key: %s", _task.TargetKey)
					return
				}
			}

			logger.Infof(ctx, "task success, task: %+v", _task)
		}); err != nil {
			return errors.Wrapf(err, "add cron job error, task: %+v", _task)
		}

		logger.Infof(ctx, "register task cron job, task: %+v", _task)
	}

	c.Start()
	return nil
}

func Init(ctx context.Context, conf *config.Config) error {
	return nil
}

func notice(ctx context.Context, _task *task.Task, err error) error {
	return email.Client().Send(ctx, mail.NewMessage().WithTitle("data-backup task error").
		WithReceiver("86744316@qq.com").
		WithDate(time.Now()).
		WithSender("data-backup-service").
		WithText(fmt.Sprintf("task: %+v\nerror: %v", _task, err)))
}
