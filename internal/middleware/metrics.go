package middleware

import (
	"strconv"
	"time"

	"git.uhomes.net/uhs-go/go-bisub/internal/pkg/metrics"
	"github.com/gin-gonic/gin"
)

// MetricsMiddleware 指标采集中间件
func MetricsMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// 增加活跃连接数
		metrics.SetActiveConnections(serviceName, 1) // 简化实现，实际应该维护计数器
		
		// 处理请求
		c.Next()
		
		// 计算耗时
		duration := time.Since(start)
		
		// 获取请求信息
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		status := c.Writer.Status()
		
		// 获取请求和响应大小
		requestSize := c.Request.ContentLength
		responseSize := int64(c.Writer.Size())
		
		// 记录指标
		metrics.RecordHTTPRequest(
			serviceName,
			method,
			path,
			status,
			duration,
			requestSize,
			responseSize,
		)
		
		// 记录错误
		if status >= 400 {
			errorType := "client_error"
			if status >= 500 {
				errorType = "server_error"
			}
			metrics.RecordError(serviceName, errorType, strconv.Itoa(status))
		}
	}
}
