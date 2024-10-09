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

	SourceId string `json:"source_id" mapstructure:"source_id"`
	TargetId string `json:"target_id" mapstructure:"target_id"`
	Source   source.Source
	Target   target.Target
}

var tasks = make([]*Task, 0)

func RegisterTasks(_ context.Context, defaultCron string, originTasks []*Task) error {
	for _, task := range originTasks {
		if task.Cron == "" {
			task.Cron = defaultCron
		}
		if task.SourceId == "" || task.TargetId == "" {
			return errors.New("task source or target is empty")
		}

		_source := source.GetSource(task.SourceId)
		if _source == nil {
			return errors.Errorf("source not found, source_id: %s", task.SourceId)
		}
		_target := target.GetTarget(task.TargetId)
		if _target == nil {
			return errors.Errorf("target not found, target_id: %s", task.TargetId)
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
