package fx

import (
	"context"
	"fmt"
	"log/slog"

	"git.uhomes.net/uhs-go/go-bisub/internal/config"
	"git.uhomes.net/uhs-go/go-bisub/internal/handler"
	"git.uhomes.net/uhs-go/go-bisub/internal/middleware"
	"git.uhomes.net/uhs-go/go-bisub/internal/models"
	"git.uhomes.net/uhs-go/go-bisub/internal/pkg/logger"
	"git.uhomes.net/uhs-go/go-bisub/internal/pkg/metrics"
	"git.uhomes.net/uhs-go/go-bisub/internal/repository"
	"git.uhomes.net/uhs-go/go-bisub/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConfigModule provides configuration
var ConfigModule = fx.Module("config",
	fx.Provide(config.Load),
)

// LoggerModule provides logger
var LoggerModule = fx.Module("logger",
	fx.Provide(func(cfg *config.Config) *logger.Logger {
		isDev := cfg.Logging.Level == "debug"
		return logger.NewLogger(cfg.Logging.Level, isDev)
	}),
	fx.Invoke(func(l *logger.Logger, cfg *config.Config) {
		logger.SetDefault(l)
		slog.Info("Logger initialized")
		
		// 初始化文件日志
		if cfg.Logging.FileLogEnabled {
			logDir := cfg.Logging.FileLogDir
			if logDir == "" {
				logDir = "./logs"
			}
			if err := logger.InitFileLogger(logDir); err != nil {
				slog.Error("Failed to initialize file logger", "error", err)
			} else {
				slog.Info("File logger initialized", "dir", logDir)
			}
		}
	}),
)

// DatabaseModule provides database connections
var DatabaseModule = fx.Module("database",
	fx.Provide(
		func(cfg *config.Config) (*gorm.DB, error) {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				cfg.Database.Primary.Username,
				cfg.Database.Primary.Password,
				cfg.Database.Primary.Host,
				cfg.Database.Primary.Port,
				cfg.Database.Primary.Database,
			)

			// 配置GORM日志
			gormConfig := &gorm.Config{}
			if cfg.Logging.FileLogEnabled {
				fileLogger := logger.GetFileLogger()
				gormLogger := logger.NewGormLogger(fileLogger)
				gormConfig.Logger = gormLogger
			}

			db, err := gorm.Open(mysql.Open(dsn), gormConfig)
			if err != nil {
				return nil, err
			}

			sqlDB, err := db.DB()
			if err != nil {
				return nil, err
			}

			sqlDB.SetMaxIdleConns(cfg.Database.Primary.MaxIdleConns)
			sqlDB.SetMaxOpenConns(cfg.Database.Primary.MaxOpenConns)
			sqlDB.SetConnMaxLifetime(cfg.Database.Primary.ConnMaxLifetime)

			// Auto migrate
			if err := db.AutoMigrate(&models.Subscription{}, &models.SubscriptionStats{}, &models.OperationLog{}); err != nil {
				return nil, err
			}

			return db, nil
		},
		func(cfg *config.Config, primaryDB *gorm.DB) map[string]*gorm.DB {
			dataSources := make(map[string]*gorm.DB)
			dataSources["primary"] = primaryDB

			// 配置GORM日志
			gormConfig := &gorm.Config{}
			if cfg.Logging.FileLogEnabled {
				fileLogger := logger.GetFileLogger()
				gormLogger := logger.NewGormLogger(fileLogger)
				gormConfig.Logger = gormLogger
			}

			for name, dbConfig := range cfg.Database.DataSources {
				dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
					dbConfig.Username,
					dbConfig.Password,
					dbConfig.Host,
					dbConfig.Port,
					dbConfig.Database,
				)

				db, err := gorm.Open(mysql.Open(dsn), gormConfig)
				if err != nil {
					slog.Error("Failed to connect to data source", "name", name, "error", err)
					continue
				}

				sqlDB, _ := db.DB()
				sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
				sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
				sqlDB.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

				dataSources[name] = db
			}

			return dataSources
		},
	),
)

// RedisModule provides Redis client
var RedisModule = fx.Module("redis",
	fx.Provide(func(cfg *config.Config) *redis.Client {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})

		// Test connection
		if err := client.Ping(context.Background()).Err(); err != nil {
			slog.Error("Failed to connect to Redis", "error", err)
			panic(err)
		}

		return client
	}),
)

// RepositoryModule provides repositories
var RepositoryModule = fx.Module("repository",
	fx.Provide(
		repository.NewSubscriptionRepository,
		repository.NewStatsRepository,
		repository.NewRefsRepository,
		repository.NewOperationLogRepository,
	),
)

// ServiceModule provides services
var ServiceModule = fx.Module("service",
	fx.Provide(
		service.NewSubscriptionService,
		service.NewRefsService,
		service.NewOperationLogService,
	),
)

// HandlerModule provides handlers
var HandlerModule = fx.Module("handler",
	fx.Provide(
		handler.NewSubscriptionHandler,
		handler.NewRefsHandler,
		handler.NewOperationLogHandler,
	),
)

// MiddlewareModule provides middlewares
var MiddlewareModule = fx.Module("middleware",
	fx.Provide(
		middleware.NewAuthMiddleware,
		func(client *redis.Client, cfg *config.Config) *middleware.RateLimiter {
			return middleware.NewRateLimiter(client, cfg.Server.RateLimit)
		},
	),
)

// HTTPModule provides HTTP server
var HTTPModule = fx.Module("http",
	fx.Provide(NewGinEngine),
	fx.Invoke(RegisterRoutes),
	fx.Invoke(InitMetrics),
)

// NewGinEngine creates a new Gin engine
func NewGinEngine(cfg *config.Config) *gin.Engine {
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	
	// 使用指标中间件
	engine.Use(middleware.MetricsMiddleware("go-bisub"))
	
	// 使用自定义日志中间件
	if cfg.Logging.FileLogEnabled {
		if cfg.Logging.LogRequestBody && cfg.Logging.LogResponseBody {
			engine.Use(middleware.LoggerMiddleware())
		} else {
			engine.Use(middleware.SimpleLoggerMiddleware())
		}
	} else {
		engine.Use(gin.Logger())
	}

	return engine
}

// InitMetrics 初始化指标系统
func InitMetrics(cfg *config.Config) {
	// 初始化指标
	metrics.Init("go-bisub")
	slog.Info("Metrics system initialized")
}

// RegisterRoutes registers all routes
func RegisterRoutes(
	engine *gin.Engine,
	cfg *config.Config,
	subscriptionHandler *handler.SubscriptionHandler,
	refsHandler *handler.RefsHandler,
	operationLogHandler *handler.OperationLogHandler,
	authMiddleware *middleware.AuthMiddleware,
	rateLimiter *middleware.RateLimiter,
) {
	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	
	// Metrics endpoint
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes (需要 JWT 认证)
	v1 := engine.Group("/v1")
	v1.Use(rateLimiter.RateLimit())
	v1.Use(authMiddleware.JWTAuth())
	{
		// Refs
		v1.GET("/refs/subscription-types", refsHandler.GetSubscriptionTypes)
		v1.GET("/refs/subscription-statuses", refsHandler.GetSubscriptionStatuses)

		// Subscriptions
		v1.GET("/subscriptions", subscriptionHandler.GetSubscriptions)
		v1.POST("/subscriptions", subscriptionHandler.CreateSubscription)
		v1.GET("/subscriptions/:key", subscriptionHandler.GetSubscription)
		v1.GET("/subscriptions/:key/versions/:version", subscriptionHandler.GetSubscription)
		v1.PUT("/subscriptions/:key/versions/:version", subscriptionHandler.UpdateSubscription)
		v1.PATCH("/subscriptions/:key/versions/:version/status", subscriptionHandler.UpdateSubscriptionStatus)
		v1.DELETE("/subscriptions/:key/versions/:version", subscriptionHandler.DeleteSubscription)

		// Execution
		v1.POST("/subscriptions/:key/execute", subscriptionHandler.ExecuteSubscription)
		v1.POST("/subscriptions/:key/versions/:version/execute", subscriptionHandler.ExecuteSubscription)

		// Stats
		v1.GET("/subscriptions/stats", subscriptionHandler.GetStats)

		// Operation logs
		v1.GET("/operation-logs", operationLogHandler.GetOperationLogs)
	}

	// Internal API for Web UI (使用 BasicAuth，与 Web UI 共享认证)
	api := engine.Group("/api")
	api.Use(authMiddleware.BasicAuth())
	{
		// Refs
		api.GET("/refs/subscription-types", refsHandler.GetSubscriptionTypes)
		api.GET("/refs/subscription-statuses", refsHandler.GetSubscriptionStatuses)

		// Subscriptions
		api.GET("/subscriptions", subscriptionHandler.GetSubscriptions)
		api.POST("/subscriptions", subscriptionHandler.CreateSubscription)
		api.GET("/subscriptions/:key", subscriptionHandler.GetSubscription)
		api.GET("/subscriptions/:key/versions/:version", subscriptionHandler.GetSubscription)
		api.PUT("/subscriptions/:key/versions/:version", subscriptionHandler.UpdateSubscription)
		api.PATCH("/subscriptions/:key/versions/:version/status", subscriptionHandler.UpdateSubscriptionStatus)
		api.DELETE("/subscriptions/:key/versions/:version", subscriptionHandler.DeleteSubscription)

		// Execution
		api.POST("/subscriptions/:key/execute", subscriptionHandler.ExecuteSubscription)
		api.POST("/subscriptions/:key/versions/:version/execute", subscriptionHandler.ExecuteSubscription)

		// Stats
		api.GET("/subscriptions/stats", subscriptionHandler.GetStats)

		// Operation logs
		api.GET("/operation-logs", operationLogHandler.GetOperationLogs)
	}

	// Web UI
	engine.LoadHTMLGlob("web/templates/*")
	webUI := engine.Group("/admin")
	webUI.Use(authMiddleware.BasicAuth())
	{
		webUI.Static("/static", "./web/static")
		webUI.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", gin.H{"title": "BI Subscription Management"})
		})
		webUI.GET("/subscriptions", func(c *gin.Context) {
			c.HTML(200, "subscriptions.html", gin.H{"title": "Subscription Management"})
		})
		webUI.GET("/stats", func(c *gin.Context) {
			c.HTML(200, "stats.html", gin.H{"title": "Statistics"})
		})
		webUI.GET("/operation-logs", func(c *gin.Context) {
			c.HTML(200, "operation_logs.html", gin.H{"title": "Operation Logs"})
		})
	}
}