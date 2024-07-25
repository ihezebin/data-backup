package cron

import (
	"context"
	"data-backup/component/source"
	"data-backup/component/target"
	"data-backup/component/task"
	"testing"

	"github.com/ihezebin/oneness/logger"
	"github.com/robfig/cron/v3"
)

func TestRunTask(t *testing.T) {
	ctx := context.Background()

	if err := source.InitMongoSources(ctx, []*source.MongoSource{
		{
			Key:         "blog",
			DSN:         "mongodb://root:root@127.0.0.1:27017/blog?authSource=admin",
			Collections: []string{"article", "tag", "draft", "comment"},
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := target.InitOSSTargets(ctx, []*target.OSSTarget{
		{
			Key: "cos",
			Dsn: "cos://xxxx:xxxx@cos.ap-chengdu.myqcloud.com/hezebin-1258606727",
			Dir: "backup",
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := task.InitTasks(ctx, "0 1 * * * *", []*task.Task{
		{
			Cron:       "0 12 17 * * *",
			Id:         "664c72b790d71012f2753739",
			SourceType: "mongo",
			SourceKey:  "blog",
			TargetType: "oss",
			TargetKey:  "cos",
		},
	}); err != nil {
		t.Fatal(err)
	}

	c := cron.New(cron.WithSeconds())

	c.Schedule(cron.Every(0), cron.FuncJob(func() {
		tasks := task.GetTasks()
		for _, _task := range tasks {
			source := _task.GetSource()
			if source == nil {
				logger.Errorf(ctx, "task source is nil, task: %+v", _task)
				return
			}

			target := _task.GetTarget()
			if target == nil {
				logger.Errorf(ctx, "task target is nil, task: %+v", _task)
				return
			}

			data, err := source.Export(ctx)
			if err != nil {
				logger.Errorf(ctx, "source export error, task: %+v, err: %+v", _task, err)
				return
			}

			for _, datum := range data {
				err = target.Store(ctx, datum.Key, datum.Content)
				if err != nil {
					logger.Errorf(ctx, "target store error, task: %+v, err: %+v", _task, err)
					return
				}
			}
		}
		c.Stop()
	}))

	c.Run()
}
