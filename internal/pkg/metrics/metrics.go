package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics 指标收集器（兼容 gox/metrics 规范）
type Metrics struct {
	// HTTP 请求指标
	RequestTotal    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	RequestSize     *prometheus.HistogramVec
	ResponseSize    *prometheus.HistogramVec
	
	// 系统指标
	ActiveConnections *prometheus.GaugeVec
	CPUUsage          *prometheus.GaugeVec
	MemoryUsage       *prometheus.GaugeVec
	DiskUsage         *prometheus.GaugeVec
	
	// 数据库指标
	DBConnections     *prometheus.GaugeVec
	DBQueryDuration   *prometheus.HistogramVec
	DBQueryTotal      *prometheus.CounterVec
	DBSlowQueryTotal  *prometheus.CounterVec
	
	// 业务指标
	SubscriptionTotal *prometheus.CounterVec
	ExecutionTotal    *prometheus.CounterVec
	ExecutionDuration *prometheus.HistogramVec
	ErrorTotal        *prometheus.CounterVec
}

var globalMetrics *Metrics

// Init 初始化指标系统
func Init(serviceName string) *Metrics {
	m := &Metrics{
		// HTTP 请求计数
		RequestTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"service", "method", "endpoint", "status"},
		),
		
		// HTTP 请求延迟
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"service", "method", "endpoint", "status"},
		),
		
		// HTTP 请求大小
		RequestSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_size_bytes",
				Help:    "HTTP request size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"service", "method", "endpoint"},
		),
		
		// HTTP 响应大小
		ResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"service", "method", "endpoint"},
		),
		
		// 活跃连接数
		ActiveConnections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_active_connections",
				Help: "Number of active HTTP connections",
			},
			[]string{"service"},
		),
		
		// CPU 使用率
		CPUUsage: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "system_cpu_usage_percent",
				Help: "CPU usage percentage",
			},
			[]string{"service"},
		),
		
		// 内存使用
		MemoryUsage: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "system_memory_usage_bytes",
				Help: "Memory usage in bytes",
			},
			[]string{"service", "type"}, // type: used, free, total
		),
		
		// 磁盘使用
		DiskUsage: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "system_disk_usage_percent",
				Help: "Disk usage percentage",
			},
			[]string{"service", "mount"},
		),
		
		// 数据库连接数
		DBConnections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_connections",
				Help: "Number of database connections",
			},
			[]string{"service", "database", "state"}, // state: idle, active, total
		),
		
		// 数据库查询延迟
		DBQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
			},
			[]string{"service", "database", "operation"},
		),
		
		// 数据库查询总数
		DBQueryTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"service", "database", "operation", "status"},
		),
		
		// 慢查询总数
		DBSlowQueryTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_slow_queries_total",
				Help: "Total number of slow database queries",
			},
			[]string{"service", "database"},
		),
		
		// 订阅总数
		SubscriptionTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "subscription_total",
				Help: "Total number of subscriptions",
			},
			[]string{"service", "type", "status"},
		),
		
		// 执行总数
		ExecutionTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "execution_total",
				Help: "Total number of executions",
			},
			[]string{"service", "subscription_key", "status"},
		),
		
		// 执行延迟
		ExecutionDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "execution_duration_seconds",
				Help:    "Execution duration in seconds",
				Buckets: []float64{.1, .25, .5, 1, 2.5, 5, 10, 30, 60, 120},
			},
			[]string{"service", "subscription_key"},
		),
		
		// 错误总数
		ErrorTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "errors_total",
				Help: "Total number of errors",
			},
			[]string{"service", "type", "code"},
		),
	}
	
	globalMetrics = m
	return m
}

// GetMetrics 获取全局指标实例
func GetMetrics() *Metrics {
	if globalMetrics == nil {
		return Init("go-bisub")
	}
	return globalMetrics
}

// RecordHTTPRequest 记录 HTTP 请求
func RecordHTTPRequest(service, method, endpoint string, status int, duration time.Duration, requestSize, responseSize int64) {
	m := GetMetrics()
	statusStr := statusCodeToString(status)
	
	// 请求计数
	m.RequestTotal.WithLabelValues(service, method, endpoint, statusStr).Inc()
	
	// 请求延迟
	m.RequestDuration.WithLabelValues(service, method, endpoint, statusStr).Observe(duration.Seconds())
	
	// 请求大小
	if requestSize > 0 {
		m.RequestSize.WithLabelValues(service, method, endpoint).Observe(float64(requestSize))
	}
	
	// 响应大小
	if responseSize > 0 {
		m.ResponseSize.WithLabelValues(service, method, endpoint).Observe(float64(responseSize))
	}
}

// RecordDBQuery 记录数据库查询
func RecordDBQuery(service, database, operation string, duration time.Duration, err error) {
	m := GetMetrics()
	
	status := "success"
	if err != nil {
		status = "error"
	}
	
	// 查询计数
	m.DBQueryTotal.WithLabelValues(service, database, operation, status).Inc()
	
	// 查询延迟
	m.DBQueryDuration.WithLabelValues(service, database, operation).Observe(duration.Seconds())
	
	// 慢查询（>200ms）
	if duration > 200*time.Millisecond {
		m.DBSlowQueryTotal.WithLabelValues(service, database).Inc()
	}
}

// RecordExecution 记录订阅执行
func RecordExecution(service, subscriptionKey string, duration time.Duration, err error) {
	m := GetMetrics()
	
	status := "success"
	if err != nil {
		status = "error"
	}
	
	// 执行计数
	m.ExecutionTotal.WithLabelValues(service, subscriptionKey, status).Inc()
	
	// 执行延迟
	m.ExecutionDuration.WithLabelValues(service, subscriptionKey).Observe(duration.Seconds())
}

// RecordError 记录错误
func RecordError(service, errorType, errorCode string) {
	m := GetMetrics()
	m.ErrorTotal.WithLabelValues(service, errorType, errorCode).Inc()
}

// SetActiveConnections 设置活跃连接数
func SetActiveConnections(service string, count int) {
	m := GetMetrics()
	m.ActiveConnections.WithLabelValues(service).Set(float64(count))
}

// SetDBConnections 设置数据库连接数
func SetDBConnections(service, database, state string, count int) {
	m := GetMetrics()
	m.DBConnections.WithLabelValues(service, database, state).Set(float64(count))
}

// SetCPUUsage 设置 CPU 使用率
func SetCPUUsage(service string, percent float64) {
	m := GetMetrics()
	m.CPUUsage.WithLabelValues(service).Set(percent)
}

// SetMemoryUsage 设置内存使用
func SetMemoryUsage(service, memType string, bytes int64) {
	m := GetMetrics()
	m.MemoryUsage.WithLabelValues(service, memType).Set(float64(bytes))
}

// SetDiskUsage 设置磁盘使用率
func SetDiskUsage(service, mount string, percent float64) {
	m := GetMetrics()
	m.DiskUsage.WithLabelValues(service, mount).Set(percent)
}

// Helper function to convert status code to string
func statusCodeToString(code int) string {
	if code >= 200 && code < 300 {
		return "2xx"
	} else if code >= 300 && code < 400 {
		return "3xx"
	} else if code >= 400 && code < 500 {
		return "4xx"
	} else if code >= 500 {
		return "5xx"
	}
	return "unknown"
}
