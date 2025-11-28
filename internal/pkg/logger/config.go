package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 日志配置
type Config struct {
	Level           string // 日志级别: debug, info, warn, error
	Format          string // 日志格式: json, console
	OutputPaths     []string
	ErrorOutputPaths []string
	Development     bool
	EnableCaller    bool
	EnableStacktrace bool
	
	// 文件日志配置
	FileLogEnabled  bool
	FileLogDir      string
	LogRequestBody  bool
	LogResponseBody bool
	
	// 日志轮转配置
	MaxSize    int  // MB
	MaxAge     int  // days
	MaxBackups int
	Compress   bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:            "info",
		Format:           "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Development:      false,
		EnableCaller:     true,
		EnableStacktrace: false,
		FileLogEnabled:   false,
		FileLogDir:       "./logs",
		LogRequestBody:   false,
		LogResponseBody:  false,
		MaxSize:          100,
		MaxAge:           30,
		MaxBackups:       10,
		Compress:         true,
	}
}

// NewLoggerFromConfig 从配置创建 logger
func NewLoggerFromConfig(cfg *Config) (*Logger, error) {
	// 解析日志级别
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置输出路径
	outputPaths := cfg.OutputPaths
	if cfg.FileLogEnabled {
		// 确保日志目录存在
		if err := os.MkdirAll(cfg.FileLogDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
		
		// 添加文件输出
		dateStr := time.Now().Format("060102")
		logFile := filepath.Join(cfg.FileLogDir, fmt.Sprintf("%s.log", dateStr))
		outputPaths = append(outputPaths, logFile)
	}

	// 创建 zap 配置
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       cfg.Development,
		DisableCaller:     !cfg.EnableCaller,
		DisableStacktrace: !cfg.EnableStacktrace,
		Encoding:          cfg.Format,
		EncoderConfig:     encoderConfig,
		OutputPaths:       outputPaths,
		ErrorOutputPaths:  cfg.ErrorOutputPaths,
	}

	// 构建 zap logger
	zapLogger, err := zapConfig.Build(
		zap.AddCallerSkip(1), // 跳过包装层
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	// 创建 slog handler
	slogHandler := NewZapHandler(zapLogger)
	slogLogger := slog.New(slogHandler)

	return &Logger{
		zap:  zapLogger,
		slog: slogLogger,
	}, nil
}

// NewDevelopmentLogger 创建开发环境 logger
func NewDevelopmentLogger() (*Logger, error) {
	cfg := DefaultConfig()
	cfg.Level = "debug"
	cfg.Format = "console"
	cfg.Development = true
	cfg.EnableStacktrace = true
	return NewLoggerFromConfig(cfg)
}

// NewProductionLogger 创建生产环境 logger
func NewProductionLogger() (*Logger, error) {
	cfg := DefaultConfig()
	cfg.Level = "info"
	cfg.Format = "json"
	cfg.Development = false
	cfg.EnableStacktrace = false
	return NewLoggerFromConfig(cfg)
}

// WithFields 添加字段到 logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	
	newZap := l.zap.With(zapFields...)
	newSlog := slog.New(NewZapHandler(newZap))
	
	return &Logger{
		zap:  newZap,
		slog: newSlog,
	}
}

// WithRequestID 添加 request_id 到 logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithFields(map[string]interface{}{
		"request_id": requestID,
	})
}

// Named 创建命名 logger
func (l *Logger) Named(name string) *Logger {
	newZap := l.zap.Named(name)
	newSlog := slog.New(NewZapHandler(newZap))
	
	return &Logger{
		zap:  newZap,
		slog: newSlog,
	}
}

// Sync 同步日志
func (l *Logger) Sync() error {
	return l.zap.Sync()
}
