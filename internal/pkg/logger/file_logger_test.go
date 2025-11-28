package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileLogger(t *testing.T) {
	// 创建临时目录
	tmpDir := filepath.Join(os.TempDir(), "test_logs")
	defer os.RemoveAll(tmpDir)

	// 初始化日志记录器
	logger, err := NewFileLogger(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}
	defer logger.Close()

	// 测试API日志
	t.Run("LogAPI", func(t *testing.T) {
		entry := &APILogEntry{
			RequestID:  "test-request-id",
			Method:     "POST",
			Path:       "/api/test",
			ClientIP:   "127.0.0.1",
			StatusCode: 200,
			Duration:   100,
		}

		err := logger.LogAPI(entry)
		if err != nil {
			t.Errorf("Failed to log API entry: %v", err)
		}

		// 验证文件是否创建
		dateStr := time.Now().Format("060102")
		apiLogPath := filepath.Join(tmpDir, dateStr+".log")
		if _, err := os.Stat(apiLogPath); os.IsNotExist(err) {
			t.Errorf("API log file not created: %s", apiLogPath)
		}
	})

	// 测试SQL日志
	t.Run("LogSQL", func(t *testing.T) {
		entry := &SQLLogEntry{
			RequestID:    "test-request-id",
			SQL:          "SELECT * FROM test",
			Duration:     50,
			RowsAffected: 10,
		}

		err := logger.LogSQL(entry)
		if err != nil {
			t.Errorf("Failed to log SQL entry: %v", err)
		}

		// 验证文件是否创建
		dateStr := time.Now().Format("060102")
		sqlLogPath := filepath.Join(tmpDir, dateStr+"_sql.log")
		if _, err := os.Stat(sqlLogPath); os.IsNotExist(err) {
			t.Errorf("SQL log file not created: %s", sqlLogPath)
		}
	})
}

func TestLogAPISimple(t *testing.T) {
	// 创建临时目录
	tmpDir := filepath.Join(os.TempDir(), "test_logs_simple")
	defer os.RemoveAll(tmpDir)

	// 初始化全局日志记录器
	err := InitFileLogger(tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize file logger: %v", err)
	}

	// 测试简化API日志
	LogAPISimple("test-id", "GET", "/api/test", "127.0.0.1", 200, 100, nil)

	// 验证文件是否创建
	dateStr := time.Now().Format("060102")
	apiLogPath := filepath.Join(tmpDir, dateStr+".log")
	if _, err := os.Stat(apiLogPath); os.IsNotExist(err) {
		t.Errorf("API log file not created: %s", apiLogPath)
	}
}

func TestLogSQLSimple(t *testing.T) {
	// 创建临时目录
	tmpDir := filepath.Join(os.TempDir(), "test_logs_sql_simple")
	defer os.RemoveAll(tmpDir)

	// 初始化全局日志记录器
	err := InitFileLogger(tmpDir)
	if err != nil {
		t.Fatalf("Failed to initialize file logger: %v", err)
	}

	// 测试简化SQL日志
	LogSQLSimple("test-id", "SELECT * FROM test", "test_db", 50, 10, nil)

	// 验证文件是否创建
	dateStr := time.Now().Format("060102")
	sqlLogPath := filepath.Join(tmpDir, dateStr+"_sql.log")
	if _, err := os.Stat(sqlLogPath); os.IsNotExist(err) {
		t.Errorf("SQL log file not created: %s", sqlLogPath)
	}
}
