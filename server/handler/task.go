package handler

import (
	"context"
	"data-backup/application/dto"
	"data-backup/application/service"

	"github.com/gin-gonic/gin"
	valication "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ihezebin/oneness/httpserver"
	"github.com/ihezebin/oneness/logger"
)

type TaskHandler struct {
	logger  *logger.Entry
	service *service.TaskApplicationService
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

func (h *TaskHandler) Init(router gin.IRouter) {
	h.logger = logger.WithField("handler", "task")
	h.service = service.NewTaskApplicationService(h.logger)

	// registry http handler
	if router != nil {
		task := router.Group("task")
		task.POST("/restore", httpserver.NewHandlerFunc(h.Restore))
		task.POST("/trigger", httpserver.NewHandlerFunc(h.Trigger))
	}

}

// Restore https://github.com/swaggo/swag/blob/master/README_zh-CN.md#api%E6%93%8D%E4%BD%9C
// @Summary 通过备份任务来恢复数据
// @Description
// @Tags task
// @Accept json
// @Produce json
// @Param req body dto.TaskRestoreReq true "任务参数"
// @Success 200 {object} server.Body{data=dto.TaskRestoreResp} "成功时如下结构；错误时 code 非 0, message 包含错误信息, 不包含 data"
// @Router /task/restore [post]
func (h *TaskHandler) Restore(ctx context.Context, req *dto.TaskRestoreReq) (*dto.TaskRestoreResp, error) {
	if err := valication.ValidateStruct(req); err != nil {
		h.logger.WithError(err).Errorf(ctx, "validate struct error, req: %v", req)
		return nil, httpserver.ErrorWithBadRequest()
	}

	return h.service.Restore(ctx, req)

}

// Trigger https://github.com/swaggo/swag/blob/master/README_zh-CN.md#api%E6%93%8D%E4%BD%9C
// @Summary 主动触发任务
// @Description
// @Tags task
// @Accept json
// @Produce json
// @Param req body dto.TaskTriggerReq true "任务参数"
// @Success 200 {object} server.Body{data=dto.TaskTriggerResp} "成功时如下结构；错误时 code 非 0, message 包含错误信息, 不包含 data"
// @Router /task/trigger [post]
func (h *TaskHandler) Trigger(ctx context.Context, req *dto.TaskTriggerReq) (*dto.TaskTriggerResp, error) {
	if err := valication.ValidateStruct(req,
		valication.Field(&req.TaskId, valication.Required),
	); err != nil {
		h.logger.WithError(err).Errorf(ctx, "validate struct error, req: %v", req)
		return nil, httpserver.ErrorWithBadRequest()
	}

	return h.service.Trigger(ctx, req)

}
