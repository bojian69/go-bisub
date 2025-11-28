package logger

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// FileLogger 文件日志记录器
// 使用 slog 作为前端接口，zap 作为后端实现
type FileLogger struct {
	logDir      string
	mu          sync.RWMutex
	apiFile     *os.File
	sqlFile     *os.File
	currentDate string
	zapLogger   *zap.Logger
	slogLogger  *slog.Logger
}

// APILogEntry API日志条目
type APILogEntry struct {
	Timestamp    string      `json:"timestamp"`
	RequestID    string      `json:"request_id"`
	Method       string      `json:"method"`
	Path         string      `json:"path"`
	ClientIP     string      `json:"client_ip"`
	UserAgent    string      `json:"user_agent"`
	StatusCode   int         `json:"status_code"`
	Duration     int64       `json:"duration_ms"`
	RequestBody  interface{} `json:"request_body,omitempty"`
	ResponseBody interface{} `json:"response_body,omitempty"`
	ErrorMessage string      `json:"error_message,omitempty"`
}

// SQLLogEntry SQL日志条目
type SQLLogEntry struct {
	Timestamp    string                 `json:"timestamp"`
	RequestID    string                 `json:"request_id"`
	SQL          string                 `json:"sql"`
	Duration     int64                  `json:"duration_ms"`
	RowsAffected int64                  `json:"rows_affected"`
	Error        string                 `json:"error,omitempty"`
	Variables    map[string]interface{} `json:"variables,omitempty"`
	Database     string                 `json:"database,omitempty"`
}

var (
	globalLogger *FileLogger
	once         sync.Once
)

// InitFileLogger 初始化文件日志记录器
func InitFileLogger(logDir string) error {
	var err error
	once.Do(func() {
		globalLogger, err = NewFileLogger(logDir)
	})
	return err
}

// GetFileLogger 获取全局文件日志记录器
func GetFileLogger() *FileLogger {
	if globalLogger == nil {
		// 如果未初始化，使用默认目录
		_ = InitFileLogger("./logs")
	}
	return globalLogger
}

// NewFileLogger 创建新的文件日志记录器
func NewFileLogger(logDir string) (*FileLogger, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	logger := &FileLogger{
		logDir:      logDir,
		currentDate: time.Now().Format("060102"), // YYMMDD格式
	}

	// 打开日志文件
	if err := logger.rotateFiles(); err != nil {
		return nil, err
	}

	// 创建 zap logger 用于结构化日志
	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{
		filepath.Join(logDir, time.Now().Format("060102")+".log"),
	}
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	
	zapLogger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create zap logger: %w", err)
	}
	
	logger.zapLogger = zapLogger
	
	// 创建 slog logger
	handler := NewZapHandler(zapLogger)
	logger.slogLogger = slog.New(handler)

	// 启动定时检查日期变化的goroutine
	go logger.watchDateChange()

	return logger, nil
}

// Slog 返回 slog.Logger 实例
func (l *FileLogger) Slog() *slog.Logger {
	return l.slogLogger
}

// Zap 返回 zap.Logger 实例
func (l *FileLogger) Zap() *zap.Logger {
	return l.zapLogger
}

// rotateFiles 轮转日志文件
func (l *FileLogger) rotateFiles() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 关闭旧文件
	if l.apiFile != nil {
		l.apiFile.Close()
	}
	if l.sqlFile != nil {
		l.sqlFile.Close()
	}

	// 打开新文件
	dateStr := time.Now().Format("060102")
	l.currentDate = dateStr

	apiPath := filepath.Join(l.logDir, fmt.Sprintf("%s.log", dateStr))
	sqlPath := filepath.Join(l.logDir, fmt.Sprintf("%s_sql.log", dateStr))

	var err error
	l.apiFile, err = os.OpenFile(apiPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open API log file: %w", err)
	}

	l.sqlFile, err = os.OpenFile(sqlPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open SQL log file: %w", err)
	}

	return nil
}

// watchDateChange 监控日期变化
func (l *FileLogger) watchDateChange() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		currentDate := time.Now().Format("060102")
		if currentDate != l.currentDate {
			if err := l.rotateFiles(); err != nil {
				fmt.Printf("Failed to rotate log files: %v\n", err)
			}
		}
	}
}

// LogAPI 记录API日志
func (l *FileLogger) LogAPI(entry *APILogEntry) error {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.apiFile == nil {
		return fmt.Errorf("API log file not initialized")
	}

	// 设置时间戳
	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().Format("2006-01-02 15:04:05.000")
	}

	// 序列化为JSON
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal API log entry: %w", err)
	}

	// 写入文件
	_, err = l.apiFile.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write API log: %w", err)
	}

	return nil
}

// LogSQL 记录SQL日志
func (l *FileLogger) LogSQL(entry *SQLLogEntry) error {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.sqlFile == nil {
		return fmt.Errorf("SQL log file not initialized")
	}

	// 设置时间戳
	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().Format("2006-01-02 15:04:05.000")
	}

	// 序列化为JSON
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal SQL log entry: %w", err)
	}

	// 写入文件
	_, err = l.sqlFile.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write SQL log: %w", err)
	}

	return nil
}

// Close 关闭日志文件
func (l *FileLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var errs []error

	if l.apiFile != nil {
		if err := l.apiFile.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if l.sqlFile != nil {
		if err := l.sqlFile.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing log files: %v", errs)
	}

	return nil
}

// LogAPISimple 简化的API日志记录
func LogAPISimple(requestID, method, path, clientIP string, statusCode int, duration int64, err error) {
	logger := GetFileLogger()
	entry := &APILogEntry{
		RequestID:  requestID,
		Method:     method,
		Path:       path,
		ClientIP:   clientIP,
		StatusCode: statusCode,
		Duration:   duration,
	}
	if err != nil {
		entry.ErrorMessage = err.Error()
	}
	_ = logger.LogAPI(entry)
}

// LogSQLSimple 简化的SQL日志记录
func LogSQLSimple(requestID, sql, database string, duration int64, rowsAffected int64, err error) {
	logger := GetFileLogger()
	entry := &SQLLogEntry{
		RequestID:    requestID,
		SQL:          sql,
		Database:     database,
		Duration:     duration,
		RowsAffected: rowsAffected,
	}
	if err != nil {
		entry.Error = err.Error()
	}
	_ = logger.LogSQL(entry)
}
