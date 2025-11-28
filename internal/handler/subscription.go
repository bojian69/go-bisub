package handler

import (
	"net/http"
	"strconv"
	"time"

	"git.uhomes.net/uhs-go/go-bisub/internal/models"
	"git.uhomes.net/uhs-go/go-bisub/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	service    *service.SubscriptionService
	logService *service.OperationLogService
}

func NewSubscriptionHandler(service *service.SubscriptionService, logService *service.OperationLogService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service:    service,
		logService: logService,
	}
}

// APIResponse 标准API响应
type APIResponse struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data,omitempty"`
	Metadata  interface{} `json:"metadata,omitempty"`
}

// CreateSubscription 创建订阅
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	startTime := time.Now()
	var req models.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 记录操作日志
		h.logOperation(c, models.OpTypeCreate, "subscription", "", models.OpStatusFailed, time.Since(startTime), err.Error(), req, nil)
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	// 从JWT中获取用户ID（简化实现）
	creatorID := uint64(1) // 实际应该从JWT token中解析

	subscription, err := h.service.CreateSubscription(c.Request.Context(), &req, creatorID)
	if err != nil {
		// 记录操作日志
		h.logOperation(c, models.OpTypeCreate, "subscription", req.SubKey, models.OpStatusFailed, time.Since(startTime), err.Error(), req, nil)
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:      "INTERNAL_ERROR",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	// 记录操作日志
	h.logOperation(c, models.OpTypeCreate, "subscription", req.SubKey, models.OpStatusSuccess, time.Since(startTime), "", req, subscription)

	c.JSON(http.StatusCreated, APIResponse{
		Code:      "OK",
		Message:   "订阅创建成功",
		RequestID: getRequestID(c),
		Data:      subscription,
	})
}

// ExecuteSubscription 执行订阅
func (h *SubscriptionHandler) ExecuteSubscription(c *gin.Context) {
	startTime := time.Now()
	subType := c.DefaultQuery("type", "A") // 默认为分析数据
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   "subscription key is required",
			RequestID: getRequestID(c),
		})
		return
	}

	var version *uint8
	if versionStr := c.Param("version"); versionStr != "" {
		if v, err := strconv.ParseUint(versionStr, 10, 8); err == nil {
			ver := uint8(v)
			version = &ver
		}
	}

	var req models.ExecuteSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	// 设置默认超时
	if req.Timeout == 0 {
		req.Timeout = 120000 // 120秒
	}

	clientIP := c.ClientIP()
	apiURL := c.Request.URL.String()

	results, err := h.service.ExecuteSubscription(c.Request.Context(), subType, key, version, &req, clientIP, apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:      "INTERNAL_ERROR",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	// 记录操作日志
	h.logOperation(c, models.OpTypeExecute, "subscription", key, models.OpStatusSuccess, time.Since(startTime), "", req, results)

	c.JSON(http.StatusOK, APIResponse{
		Code:      "OK",
		Message:   "执行成功",
		RequestID: getRequestID(c),
		Data:      results,
	})
}

// GetSubscriptions 获取订阅列表
func (h *SubscriptionHandler) GetSubscriptions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	// 获取搜索参数
	subKey := c.Query("sub_key")
	title := c.Query("title")
	status := c.Query("status")

	subscriptions, total, err := h.service.GetSubscriptions(c.Request.Context(), limit, offset, subKey, title, status)
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
		Data: map[string]interface{}{
			"items": subscriptions,
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

// GetSubscription 获取订阅详情
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	subType := c.DefaultQuery("type", "A") // 默认为分析数据
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   "subscription key is required",
			RequestID: getRequestID(c),
		})
		return
	}

	var version *uint8
	if versionStr := c.Param("version"); versionStr != "" {
		if v, err := strconv.ParseUint(versionStr, 10, 8); err == nil {
			ver := uint8(v)
			version = &ver
		}
	}

	subscription, err := h.service.GetSubscription(c.Request.Context(), subType, key, version)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:      "NOT_FOUND",
			Message:   "订阅不存在",
			RequestID: getRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:      "OK",
		Message:   "获取成功",
		RequestID: getRequestID(c),
		Data:      subscription,
	})
}

// UpdateSubscription 更新订阅
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	subType := c.DefaultQuery("type", "A")
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   "subscription key is required",
			RequestID: getRequestID(c),
		})
		return
	}

	var version uint8
	if versionStr := c.Param("version"); versionStr != "" {
		if v, err := strconv.ParseUint(versionStr, 10, 8); err == nil {
			version = uint8(v)
		} else {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:      "INVALID_PARAMETER",
				Message:   "invalid version",
				RequestID: getRequestID(c),
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   "version is required",
			RequestID: getRequestID(c),
		})
		return
	}

	var req models.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	subscription, err := h.service.UpdateSubscription(c.Request.Context(), subType, key, version, &req)
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
		Message:   "更新成功",
		RequestID: getRequestID(c),
		Data:      subscription,
	})
}

// UpdateSubscriptionStatus 更新订阅状态
func (h *SubscriptionHandler) UpdateSubscriptionStatus(c *gin.Context) {
	subType := c.DefaultQuery("type", "A")
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   "subscription key is required",
			RequestID: getRequestID(c),
		})
		return
	}

	var version uint8
	if versionStr := c.Param("version"); versionStr != "" {
		if v, err := strconv.ParseUint(versionStr, 10, 8); err == nil {
			version = uint8(v)
		} else {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:      "INVALID_PARAMETER",
				Message:   "invalid version",
				RequestID: getRequestID(c),
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   "version is required",
			RequestID: getRequestID(c),
		})
		return
	}

	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	if err := h.service.UpdateStatus(c.Request.Context(), subType, key, version, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:      "INTERNAL_ERROR",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:      "OK",
		Message:   "状态更新成功",
		RequestID: getRequestID(c),
	})
}

// DeleteSubscription 删除订阅
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	subType := c.DefaultQuery("type", "A")
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   "subscription key is required",
			RequestID: getRequestID(c),
		})
		return
	}

	var version uint8
	if versionStr := c.Param("version"); versionStr != "" {
		if v, err := strconv.ParseUint(versionStr, 10, 8); err == nil {
			version = uint8(v)
		} else {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:      "INVALID_PARAMETER",
				Message:   "invalid version",
				RequestID: getRequestID(c),
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   "version is required",
			RequestID: getRequestID(c),
		})
		return
	}

	if err := h.service.DeleteSubscription(c.Request.Context(), subType, key, version); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:      "INTERNAL_ERROR",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:      "OK",
		Message:   "删除成功",
		RequestID: getRequestID(c),
	})
}

// GetStats 获取统计数据
func (h *SubscriptionHandler) GetStats(c *gin.Context) {
	var req models.StatsQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:      "INVALID_PARAMETER",
			Message:   err.Error(),
			RequestID: getRequestID(c),
		})
		return
	}

	stats, err := h.service.GetStats(c.Request.Context(), &req)
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
		Data:      stats,
	})
}

func getRequestID(c *gin.Context) string {
	if requestID := c.GetHeader("X-Request-Id"); requestID != "" {
		return requestID
	}
	return uuid.New().String()
}

// logOperation 记录操作日志
func (h *SubscriptionHandler) logOperation(c *gin.Context, operation, resource, resourceID, status string, duration time.Duration, errorMsg string, requestData, responseData interface{}) {
	if h.logService == nil {
		return
	}

	userID := uint64(1) // 从上下文获取
	username := "admin" // 从上下文获取

	log := h.logService.CreateOperationLog(
		userID,
		username,
		operation,
		resource,
		resourceID,
		status,
		c.ClientIP(),
		c.GetHeader("User-Agent"),
		c.Request.URL.String(),
		c.Request.Method,
		uint32(duration.Milliseconds()),
		errorMsg,
		requestData,
		responseData,
	)

	h.logService.LogOperation(c.Request.Context(), log)
}
