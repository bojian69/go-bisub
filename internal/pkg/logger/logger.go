package logger

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger with slog interface
type Logger struct {
	zap  *zap.Logger
	slog *slog.Logger
}

// NewLogger creates a new logger instance
func NewLogger(level string, isDev bool) *Logger {
	var zapConfig zap.Config
	if isDev {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	// Set log level
	switch level {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	zapLogger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	// Create slog handler from zap
	slogHandler := NewZapHandler(zapLogger)
	slogLogger := slog.New(slogHandler)

	return &Logger{
		zap:  zapLogger,
		slog: slogLogger,
	}
}

// Zap returns the underlying zap logger
func (l *Logger) Zap() *zap.Logger {
	return l.zap
}

// Slog returns the slog logger
func (l *Logger) Slog() *slog.Logger {
	return l.slog
}

// ZapHandler implements slog.Handler using zap
type ZapHandler struct {
	zap *zap.Logger
}

func NewZapHandler(zap *zap.Logger) *ZapHandler {
	return &ZapHandler{zap: zap}
}

func (h *ZapHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.zap.Core().Enabled(zapcore.Level(level))
}

func (h *ZapHandler) Handle(ctx context.Context, record slog.Record) error {
	fields := make([]zap.Field, 0, record.NumAttrs())
	record.Attrs(func(attr slog.Attr) bool {
		fields = append(fields, zap.Any(attr.Key, attr.Value.Any()))
		return true
	})

	switch record.Level {
	case slog.LevelDebug:
		h.zap.Debug(record.Message, fields...)
	case slog.LevelInfo:
		h.zap.Info(record.Message, fields...)
	case slog.LevelWarn:
		h.zap.Warn(record.Message, fields...)
	case slog.LevelError:
		h.zap.Error(record.Message, fields...)
	}
	return nil
}

func (h *ZapHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := make([]zap.Field, len(attrs))
	for i, attr := range attrs {
		fields[i] = zap.Any(attr.Key, attr.Value.Any())
	}
	return &ZapHandler{zap: h.zap.With(fields...)}
}

func (h *ZapHandler) WithGroup(name string) slog.Handler {
	return &ZapHandler{zap: h.zap.Named(name)}
}

// SetDefault sets the default slog logger
func SetDefault(logger *Logger) {
	slog.SetDefault(logger.Slog())
}

// Default returns a default logger
func Default() *Logger {
	return NewLogger("info", false)
}