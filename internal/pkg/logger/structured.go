package logger

import (
	"context"
	"log/slog"
	"time"

	"go.uber.org/zap"
)

// StructuredLogger 结构化日志记录器
// 兼容 gox/log 接口规范，使用 slog + zap 实现
type StructuredLogger struct {
	logger *Logger
	fields map[string]interface{}
}

// NewStructuredLogger 创建结构化日志记录器
func NewStructuredLogger(logger *Logger) *StructuredLogger {
	return &StructuredLogger{
		logger: logger,
		fields: make(map[string]interface{}),
	}
}

// WithFields 添加字段
func (l *StructuredLogger) WithFields(fields map[string]interface{}) *StructuredLogger {
	newFields := make(map[string]interface{})
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	
	return &StructuredLogger{
		logger: l.logger,
		fields: newFields,
	}
}

// WithField 添加单个字段
func (l *StructuredLogger) WithField(key string, value interface{}) *StructuredLogger {
	return l.WithFields(map[string]interface{}{key: value})
}

// WithContext 从 context 中提取字段
func (l *StructuredLogger) WithContext(ctx context.Context) *StructuredLogger {
	fields := make(map[string]interface{})
	
	// 提取 trace_id
	if traceID := ctx.Value("trace_id"); traceID != nil {
		fields["trace_id"] = traceID
	}
	
	// 提取 span_id
	if spanID := ctx.Value("span_id"); spanID != nil {
		fields["span_id"] = spanID
	}
	
	// 提取 request_id
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields["request_id"] = requestID
	}
	
	// 提取 user_id
	if userID := ctx.Value("user_id"); userID != nil {
		fields["user_id"] = userID
	}
	
	return l.WithFields(fields)
}

// Debug 记录 DEBUG 级别日志
func (l *StructuredLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.log(slog.LevelDebug, msg, keysAndValues...)
}

// Info 记录 INFO 级别日志
func (l *StructuredLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log(slog.LevelInfo, msg, keysAndValues...)
}

// Warn 记录 WARN 级别日志
func (l *StructuredLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.log(slog.LevelWarn, msg, keysAndValues...)
}

// Error 记录 ERROR 级别日志
func (l *StructuredLogger) Error(msg string, keysAndValues ...interface{}) {
	l.log(slog.LevelError, msg, keysAndValues...)
}

// Fatal 记录 FATAL 级别日志并退出
func (l *StructuredLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.log(slog.LevelError, msg, keysAndValues...)
	l.logger.Sync()
	panic(msg) // 或者 os.Exit(1)
}

// log 内部日志记录方法
func (l *StructuredLogger) log(level slog.Level, msg string, keysAndValues ...interface{}) {
	// 构建 slog 属性
	attrs := make([]slog.Attr, 0, len(l.fields)+len(keysAndValues)/2)
	
	// 添加预设字段
	for k, v := range l.fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	
	// 添加动态字段
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key, ok := keysAndValues[i].(string)
			if ok {
				attrs = append(attrs, slog.Any(key, keysAndValues[i+1]))
			}
		}
	}
	
	// 使用 slog 记录
	l.logger.Slog().LogAttrs(context.Background(), level, msg, attrs...)
}

// DebugContext 带上下文的 DEBUG 日志
func (l *StructuredLogger) DebugContext(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.WithContext(ctx).Debug(msg, keysAndValues...)
}

// InfoContext 带上下文的 INFO 日志
func (l *StructuredLogger) InfoContext(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.WithContext(ctx).Info(msg, keysAndValues...)
}

// WarnContext 带上下文的 WARN 日志
func (l *StructuredLogger) WarnContext(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.WithContext(ctx).Warn(msg, keysAndValues...)
}

// ErrorContext 带上下文的 ERROR 日志
func (l *StructuredLogger) ErrorContext(ctx context.Context, msg string, keysAndValues ...interface{}) {
	l.WithContext(ctx).Error(msg, keysAndValues...)
}

// 全局结构化日志实例
var globalStructuredLogger *StructuredLogger

// InitStructuredLogger 初始化全局结构化日志
func InitStructuredLogger(logger *Logger) {
	globalStructuredLogger = NewStructuredLogger(logger)
}

// GetStructuredLogger 获取全局结构化日志
func GetStructuredLogger() *StructuredLogger {
	if globalStructuredLogger == nil {
		globalStructuredLogger = NewStructuredLogger(Default())
	}
	return globalStructuredLogger
}

// 便捷函数

// Debug 全局 DEBUG 日志
func Debug(msg string, keysAndValues ...interface{}) {
	GetStructuredLogger().Debug(msg, keysAndValues...)
}

// Info 全局 INFO 日志
func Info(msg string, keysAndValues ...interface{}) {
	GetStructuredLogger().Info(msg, keysAndValues...)
}

// Warn 全局 WARN 日志
func Warn(msg string, keysAndValues ...interface{}) {
	GetStructuredLogger().Warn(msg, keysAndValues...)
}

// Error 全局 ERROR 日志
func Error(msg string, keysAndValues ...interface{}) {
	GetStructuredLogger().Error(msg, keysAndValues...)
}

// Fatal 全局 FATAL 日志
func Fatal(msg string, keysAndValues ...interface{}) {
	GetStructuredLogger().Fatal(msg, keysAndValues...)
}

// WithFields 创建带字段的日志记录器
func WithFields(fields map[string]interface{}) *StructuredLogger {
	return GetStructuredLogger().WithFields(fields)
}

// WithField 创建带单个字段的日志记录器
func WithField(key string, value interface{}) *StructuredLogger {
	return GetStructuredLogger().WithField(key, value)
}

// WithContext 从上下文创建日志记录器
func WithContext(ctx context.Context) *StructuredLogger {
	return GetStructuredLogger().WithContext(ctx)
}

// LogEntry 日志条目（兼容 gox/log 格式）
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Service   string                 `json:"service"`
	TraceID   string                 `json:"trace_id,omitempty"`
	SpanID    string                 `json:"span_id,omitempty"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:",inline"`
}

// ToZapFields 转换为 zap 字段
func (e *LogEntry) ToZapFields() []zap.Field {
	fields := []zap.Field{
		zap.String("service", e.Service),
		zap.String("message", e.Message),
	}
	
	if e.TraceID != "" {
		fields = append(fields, zap.String("trace_id", e.TraceID))
	}
	
	if e.SpanID != "" {
		fields = append(fields, zap.String("span_id", e.SpanID))
	}
	
	for k, v := range e.Fields {
		fields = append(fields, zap.Any(k, v))
	}
	
	return fields
}

// NewLogEntry 创建日志条目
func NewLogEntry(service, message string) *LogEntry {
	return &LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   service,
		Message:   message,
		Fields:    make(map[string]interface{}),
	}
}

// WithTraceID 添加 trace_id
func (e *LogEntry) WithTraceID(traceID string) *LogEntry {
	e.TraceID = traceID
	return e
}

// WithSpanID 添加 span_id
func (e *LogEntry) WithSpanID(spanID string) *LogEntry {
	e.SpanID = spanID
	return e
}

// WithField 添加字段
func (e *LogEntry) WithField(key string, value interface{}) *LogEntry {
	e.Fields[key] = value
	return e
}

// WithFields 添加多个字段
func (e *LogEntry) WithFields(fields map[string]interface{}) *LogEntry {
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}
