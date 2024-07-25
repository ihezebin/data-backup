package service

import (
	"context"
	"data-backup/application/dto"
	"data-backup/component/source"
	"data-backup/component/target"
	"data-backup/component/task"
	"testing"

	"github.com/ihezebin/oneness/logger"
)

func TestRestoreByTask(t *testing.T) {
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
			Dsn: "cos://xxx:xxx@cos.ap-chengdu.myqcloud.com/hezebin-1258606727",
			Dir: "backup",
		},
	}); err != nil {
		t.Fatal(err)
	}

	if err := task.InitTasks(ctx, "0 1 * * * *", []*task.Task{
		{
			Id:         "664c72b790d71012f2753739",
			SourceType: "mongo",
			SourceKey:  "blog",
			TargetType: "oss",
			TargetKey:  "cos",
		},
	}); err != nil {
		t.Fatal(err)
	}

	service := NewRestoreApplicationService(logger.WithField("application", "restore"))

	resp, err := service.RestoreByTask(ctx, &dto.RestoreByTaskReq{
		TaskId: "664c72b790d71012f2753739",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", resp.Task)
}
