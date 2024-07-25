package dto

import "data-backup/component/task"

type RestoreByTaskReq struct {
	TaskId string `json:"task_id"`
}

type RestoreByTaskResp struct {
	Task *task.Task `json:"task"`
}
