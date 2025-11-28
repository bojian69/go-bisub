package middleware

import (
	"bytes"
	"io"
	"time"

	"git.uhomes.net/uhs-go/go-bisub/internal/models"
	"git.uhomes.net/uhs-go/go-bisub/internal/service"
	"github.com/gin-gonic/gin"
)

// OperationLogMiddleware 操作日志中间件
func OperationLogMiddleware(logService *service.OperationLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装响应写入器
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = w

		// 处理请求
		c.Next()

		// 计算执行时间
		duration := uint32(time.Since(startTime).Milliseconds())

		// 获取用户信息（从JWT或session中获取）
		userID := getUserID(c)
		username := getUsername(c)

		// 确定操作类型
		operation := getOperationType(c.Request.Method, c.FullPath())

		// 确定资源
		resource := getResource(c.FullPath())
		resourceID := getResourceID(c)

		// 确定状态
		status := models.OpStatusSuccess
		errorMsg := ""
		if c.Writer.Status() >= 400 {
			status = models.OpStatusFailed
			errorMsg = w.body.String()
		}

		// 创建操作日志
		log := logService.CreateOperationLog(
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
			duration,
			errorMsg,
			string(requestBody),
			w.body.String(),
		)

		// 异步记录日志
		logService.LogOperation(c.Request.Context(), log)
	}
}

func getUserID(c *gin.Context) uint64 {
	// 从JWT token或session中获取用户ID
	// 这里简化实现，实际应该从认证中间件设置的上下文中获取
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint64); ok {
			return id
		}
	}
	return 1 // 默认用户ID
}

func getUsername(c *gin.Context) string {
	// 从JWT token或session中获取用户名
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name
		}
	}
	return "admin" // 默认用户名
}

func getOperationType(method, path string) string {
	switch method {
	case "POST":
		if contains(path, "execute") {
			return models.OpTypeExecute
		}
		return models.OpTypeCreate
	case "PUT", "PATCH":
		return models.OpTypeUpdate
	case "DELETE":
		return models.OpTypeDelete
	case "GET":
		return models.OpTypeQuery
	default:
		return "UNKNOWN"
	}
}

func getResource(path string) string {
	if contains(path, "subscriptions") {
		return "subscription"
	}
	if contains(path, "stats") {
		return "stats"
	}
	if contains(path, "logs") {
		return "operation_log"
	}
	return "unknown"
}

func getResourceID(c *gin.Context) string {
	if key := c.Param("key"); key != "" {
		if version := c.Param("version"); version != "" {
			return key + ":" + version
		}
		return key
	}
	return ""
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)+1] == substr+"/" || s[len(s)-len(substr)-1:] == "/"+substr)))
}
