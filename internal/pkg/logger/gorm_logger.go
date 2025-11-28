package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger GORM日志记录器
// 使用 slog 作为前端接口，zap 作为后端实现
type GormLogger struct {
	fileLogger                *FileLogger
	slogLogger                *slog.Logger
	zapLogger                 *zap.Logger
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// NewGormLogger 创建GORM日志记录器
func NewGormLogger(fileLogger *FileLogger) *GormLogger {
	var slogLogger *slog.Logger
	var zapLogger *zap.Logger
	
	if fileLogger != nil {
		slogLogger = fileLogger.Slog()
		zapLogger = fileLogger.Zap()
	} else {
		// 使用默认 slog
		slogLogger = slog.Default()
	}
	
	return &GormLogger{
		fileLogger:                fileLogger,
		slogLogger:                slogLogger,
		zapLogger:                 zapLogger,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}
}

// LogMode 实现gorm logger接口
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info 实现gorm logger接口
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	// 可以选择记录INFO级别日志
}

// Warn 实现gorm logger接口
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	// 可以选择记录WARN级别日志
}

// Error 实现gorm logger接口
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	// 可以选择记录ERROR级别日志
}

// Trace 实现gorm logger接口 - 记录SQL执行
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// 从context中获取requestID
	requestID := getRequestIDFromContext(ctx)

	// 判断是否为慢查询
	isSlow := elapsed > l.SlowThreshold
	
	// 使用 slog 记录
	logLevel := slog.LevelInfo
	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {
		logLevel = slog.LevelError
	} else if isSlow {
		logLevel = slog.LevelWarn
	}
	
	logAttrs := []slog.Attr{
		slog.String("request_id", requestID),
		slog.String("sql", sql),
		slog.Int64("duration_ms", elapsed.Milliseconds()),
		slog.Int64("rows_affected", rows),
	}
	
	if isSlow {
		logAttrs = append(logAttrs, slog.Bool("slow_query", true))
	}
	
	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {
		logAttrs = append(logAttrs, slog.String("error", err.Error()))
	}
	
	// 使用 slog 记录
	if l.slogLogger != nil {
		l.slogLogger.LogAttrs(ctx, logLevel, "SQL Query", logAttrs...)
	}
	
	// 使用 zap 记录（高性能）
	if l.zapLogger != nil {
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("sql", sql),
			zap.Int64("duration_ms", elapsed.Milliseconds()),
			zap.Int64("rows_affected", rows),
		}
		
		if isSlow {
			fields = append(fields, zap.Bool("slow_query", true))
		}
		
		if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {
			fields = append(fields, zap.Error(err))
			l.zapLogger.Error("SQL Query", fields...)
		} else if isSlow {
			l.zapLogger.Warn("SQL Query", fields...)
		} else {
			l.zapLogger.Info("SQL Query", fields...)
		}
	}

	// 同时记录到文件（保持兼容性）
	if l.fileLogger != nil {
		entry := &SQLLogEntry{
			RequestID:    requestID,
			SQL:          sql,
			Duration:     elapsed.Milliseconds(),
			RowsAffected: rows,
		}

		if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {
			entry.Error = err.Error()
		}

		if isSlow {
			entry.SQL = fmt.Sprintf("[SLOW QUERY] %s", sql)
		}

		_ = l.fileLogger.LogSQL(entry)
	}
}

// getRequestIDFromContext 从context中获取requestID
func getRequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	
	// 尝试从context中获取requestID
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	
	return ""
}

// SetRequestID 设置requestID到context
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, "request_id", requestID)
}
