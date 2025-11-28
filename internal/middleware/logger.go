package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"git.uhomes.net/uhs-go/go-bisub/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// responseWriter 包装gin.ResponseWriter以捕获响应体
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware API日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 生成或获取requestID
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-Id", requestID)

		// 将requestID设置到context中，供GORM使用
		ctx := logger.SetRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		// 读取请求体
		var requestBody interface{}
		if c.Request.Body != nil && c.Request.Method != "GET" {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			
			// 尝试解析JSON
			if len(bodyBytes) > 0 {
				var jsonBody map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
					requestBody = jsonBody
				} else {
					requestBody = string(bodyBytes)
				}
			}
		}

		// 包装ResponseWriter以捕获响应
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime)

		// 解析响应体
		var responseBody interface{}
		if blw.body.Len() > 0 {
			var jsonResp map[string]interface{}
			if err := json.Unmarshal(blw.body.Bytes(), &jsonResp); err == nil {
				responseBody = jsonResp
			} else {
				// 如果不是JSON，只记录前1000个字符
				respStr := blw.body.String()
				if len(respStr) > 1000 {
					respStr = respStr[:1000] + "..."
				}
				responseBody = respStr
			}
		}

		// 使用 slog 记录API日志
		logAttrs := []slog.Attr{
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.Int("status_code", c.Writer.Status()),
			slog.Int64("duration_ms", duration.Milliseconds()),
		}

		// 添加请求体（如果配置启用）
		if requestBody != nil {
			logAttrs = append(logAttrs, slog.Any("request_body", requestBody))
		}

		// 添加响应体（如果配置启用）
		if responseBody != nil {
			logAttrs = append(logAttrs, slog.Any("response_body", responseBody))
		}

		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			logAttrs = append(logAttrs, slog.String("error", c.Errors.String()))
		}

		// 异步记录日志，避免阻塞请求
		go func() {
			fileLogger := logger.GetFileLogger()
			if fileLogger != nil && fileLogger.Slog() != nil {
				// 使用 slog 记录
				fileLogger.Slog().LogAttrs(c.Request.Context(), slog.LevelInfo, "API Request", logAttrs...)
			}
			
			// 同时记录到文件（保持兼容性）
			entry := &logger.APILogEntry{
				RequestID:    requestID,
				Method:       c.Request.Method,
				Path:         c.Request.URL.Path,
				ClientIP:     c.ClientIP(),
				UserAgent:    c.Request.UserAgent(),
				StatusCode:   c.Writer.Status(),
				Duration:     duration.Milliseconds(),
				RequestBody:  requestBody,
				ResponseBody: responseBody,
			}
			if len(c.Errors) > 0 {
				entry.ErrorMessage = c.Errors.String()
			}
			_ = fileLogger.LogAPI(entry)
		}()
	}
}

// SimpleLoggerMiddleware 简化版日志中间件（不记录请求/响应体）
func SimpleLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 生成或获取requestID
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-Id", requestID)

		// 将requestID设置到context中
		ctx := logger.SetRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime)

		// 使用 slog 记录简化日志
		go func() {
			slog.InfoContext(c.Request.Context(), "API Request",
				slog.String("request_id", requestID),
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
				slog.String("client_ip", c.ClientIP()),
				slog.Int("status_code", c.Writer.Status()),
				slog.Int64("duration_ms", duration.Milliseconds()),
			)
			
			// 同时记录到文件
			logger.LogAPISimple(
				requestID,
				c.Request.Method,
				c.Request.URL.Path,
				c.ClientIP(),
				c.Writer.Status(),
				duration.Milliseconds(),
				nil,
			)
		}()
	}
}
