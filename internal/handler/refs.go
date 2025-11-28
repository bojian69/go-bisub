package handler

import (
	"net/http"

	"git.uhomes.net/uhs-go/go-bisub/internal/service"
	"github.com/gin-gonic/gin"
)

type RefsHandler struct {
	service *service.RefsService
}

func NewRefsHandler(service *service.RefsService) *RefsHandler {
	return &RefsHandler{service: service}
}

// GetSubscriptionTypes 获取订阅类型列表
func (h *RefsHandler) GetSubscriptionTypes(c *gin.Context) {
	types, err := h.service.GetSubscriptionTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:      "INTERNAL_ERROR",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:      "OK",
		Message:   "获取成功",
		RequestID: getRequestID(c),
		Data:      types,
	})
}

// GetSubscriptionStatuses 获取订阅状态列表
func (h *RefsHandler) GetSubscriptionStatuses(c *gin.Context) {
	statuses, err := h.service.GetSubscriptionStatuses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:      "INTERNAL_ERROR",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:      "OK",
		Message:   "获取成功",
		RequestID: getRequestID(c),
		Data:      statuses,
	})
}
