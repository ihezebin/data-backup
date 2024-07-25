package task

import (
	"context"
	"data-backup/component/source"
	"data-backup/component/target"

	"github.com/pkg/errors"
)

type Task struct {
	Id   string `json:"id" mapstructure:"id"`
	Cron string `json:"cron" mapstructure:"cron"`

	SourceType string `json:"source_type" mapstructure:"source_type"`
	SourceKey  string `json:"source_key" mapstructure:"source_key"`
	TargetType string `json:"target_type" mapstructure:"target_type"`
	TargetKey  string `json:"target_key" mapstructure:"target_key"`
}

var tasks = make([]*Task, 0)

func InitTasks(ctx context.Context, defaultCron string, originTasks []*Task) error {
	for _, task := range originTasks {
		if task.Cron == "" {
			task.Cron = defaultCron
		}
		if task.SourceType == "" || task.SourceKey == "" || task.TargetType == "" || task.TargetKey == "" {
			return errors.New("task cron, source or target is empty")
		}

		_, ok := source.SupportSourceTypes[task.SourceType]
		if !ok {
			return errors.Errorf("task source type %s not support", task.SourceType)
		}

		_, ok = target.SupportTargetTypes[task.TargetType]
		if !ok {
			return errors.Errorf("task target type %s not support", task.TargetType)
		}

		tasks = append(tasks, task)
	}
	return nil
}

func GetTasks() []*Task {
	return tasks
}

func (t *Task) GetSource() source.Source {
	switch t.SourceType {
	case source.TypeMongoDB:
		return source.GetMongoSource(t.SourceKey)
	case source.TypeMysql:
		return source.GetMysqlSource(t.SourceKey)
	default:
		return nil
	}
}

func (t *Task) GetTarget() target.Target {
	switch t.TargetType {
	case target.TypeOss:
		return target.GetOSSTarget(t.TargetKey)
	default:
		return nil
	}
}
