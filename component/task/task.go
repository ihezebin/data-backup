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

	SourceKey string `json:"source_key" mapstructure:"source_key"`
	TargetKey string `json:"target_key" mapstructure:"target_key"`
	Source    source.Source
	Target    target.Target
}

var tasks = make([]*Task, 0)

func RegisterTasks(_ context.Context, defaultCron string, originTasks []*Task) error {
	for _, task := range originTasks {
		if task.Cron == "" {
			task.Cron = defaultCron
		}
		if task.SourceKey == "" || task.TargetKey == "" {
			return errors.New("task source or target is empty")
		}

		_source := source.GetSource(task.SourceKey)
		if _source == nil {
			return errors.Errorf("source not found, source_key: %s", task.SourceKey)
		}
		_target := target.GetTarget(task.TargetKey)
		if _target == nil {
			return errors.Errorf("target not found, target_key: %s", task.TargetKey)
		}

		task.Source = _source
		task.Target = _target

		tasks = append(tasks, task)
	}
	return nil
}

func GetTasks() []*Task {
	return tasks
}
