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

type RestoreHandler struct {
	logger  *logger.Entry
	service *service.RestoreApplicationService
}

func NewRestoreHandler() *RestoreHandler {
	return &RestoreHandler{}
}

func (h *RestoreHandler) Init(router gin.IRouter) {
	h.logger = logger.WithField("handler", "restore")
	h.service = service.NewRestoreApplicationService(h.logger)

	// registry http handler
	if router != nil {
		restore := router.Group("restore")
		restore.POST("/task", httpserver.NewHandlerFunc(h.RestoreByTask))
	}

}

// RestoreByTask https://github.com/swaggo/swag/blob/master/README_zh-CN.md#api%E6%93%8D%E4%BD%9C
// @Summary 通过备份任务来恢复数据
// @Description
// @Tags restore
// @Accept json
// @Produce json
// @Param req body dto.RestoreByTaskReq true "任务参数"
// @Success 200 {object} server.Body{data=dto.RestoreByTaskResp} "成功时如下结构；错误时 code 非 0, message 包含错误信息, 不包含 data"
// @Router /restore/login [post]
func (h *RestoreHandler) RestoreByTask(ctx context.Context, req *dto.RestoreByTaskReq) (*dto.RestoreByTaskResp, error) {
	if err := valication.ValidateStruct(req,
		valication.Field(&req.TaskId, valication.Required),
	); err != nil {
		h.logger.WithError(err).Errorf(ctx, "validate struct error, req: %v", req)
		return nil, httpserver.ErrorWithBadRequest()
	}

	return h.service.RestoreByTask(ctx, req)

}
