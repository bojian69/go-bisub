package handler

import (
	"net/http"

	"git.uhomes.net/uhs-go/go-bisub/internal/models"
	"git.uhomes.net/uhs-go/go-bisub/internal/service"
	"github.com/gin-gonic/gin"
)

type OperationLogHandler struct {
	service *service.OperationLogService
}

func NewOperationLogHandler(service *service.OperationLogService) *OperationLogHandler {
	return &OperationLogHandler{service: service}
}

// GetOperationLogs 获取操作日志列表
func (h *OperationLogHandler) GetOperationLogs(c *gin.Context) {
	var req models.OperationLogRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	logs, total, err := h.service.GetOperationLogs(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:      "INTERNAL_ERROR",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:      "OK",
		Message:   "获取成功",
		RequestID: getRequestID(c),
		Data: map[string]interface{}{
			"items": logs,
			"pagination": map[string]interface{}{
				"total":        total,
				"limit":        limit,
				"offset":       offset,
				"current_page": offset/limit + 1,
				"total_pages":  (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}
