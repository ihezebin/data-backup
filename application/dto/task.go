package dto

import "data-backup/component/task"

type TaskRestoreReq struct {
	TaskId  string   `json:"task_id"`
	TaskIds []string `json:"task_ids"`
}

type TaskRestoreResp struct {
	Tasks []*task.Task `json:"tasks"`
}

type TaskTriggerReq struct {
	TaskId string `json:"task_id"`
}

type TaskTriggerResp struct {
	Task *task.Task `json:"task"`
}
