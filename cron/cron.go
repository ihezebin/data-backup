package cron

import (
	"context"
	"data-backup/component/email"
	"data-backup/component/task"
	"data-backup/config"
	"encoding/json"
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
		taskData, _ := json.Marshal(_task)
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

			logger.Infof(ctx, "task start, task: %s", taskData)

			err = _task.Source.Backup(ctx, _task.Target)
			if err != nil {
				err = errors.Wrapf(err, "source backup error, task: %s", taskData)
				return
			}

			logger.Infof(ctx, "task success, task: %s", taskData)
		}); err != nil {
			return errors.Wrapf(err, "add cron job error, task: %+v", _task)
		}

		logger.Infof(ctx, "register task cron job, task: %s", taskData)
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
