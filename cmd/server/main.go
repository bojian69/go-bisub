package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.uhomes.net/uhs-go/go-bisub/internal/config"
	"git.uhomes.net/uhs-go/go-bisub/internal/handler"
	"git.uhomes.net/uhs-go/go-bisub/internal/middleware"
	"git.uhomes.net/uhs-go/go-bisub/internal/models"
	"git.uhomes.net/uhs-go/go-bisub/internal/repository"
	"git.uhomes.net/uhs-go/go-bisub/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	primaryDB, err := initDB(&cfg.Database.Primary)
	if err != nil {
		log.Fatalf("Failed to connect to primary database: %v", err)
	}

	// 自动迁移
	if err := primaryDB.AutoMigrate(&models.Subscription{}, &models.SubscriptionStats{}, &models.OperationLog{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化数据源
	dataSources := make(map[string]*gorm.DB)
	for name, dbConfig := range cfg.Database.DataSources {
		db, err := initDB(&dbConfig)
		if err != nil {
			log.Printf("Failed to connect to data source %s: %v", name, err)
			continue
		}
		dataSources[name] = db
	}

	// 初始化Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试Redis连接
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 初始化仓库
	subscriptionRepo := repository.NewSubscriptionRepository(primaryDB)
	statsRepo := repository.NewStatsRepository(primaryDB)
	refsRepo := repository.NewRefsRepository(primaryDB)
	operationLogRepo := repository.NewOperationLogRepository(primaryDB)

	// 初始化服务
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, statsRepo, dataSources, cfg)
	refsService := service.NewRefsService(refsRepo)
	operationLogService := service.NewOperationLogService(operationLogRepo)

	// 初始化处理器
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService, operationLogService)
	refsHandler := handler.NewRefsHandler(refsService)
	operationLogHandler := handler.NewOperationLogHandler(operationLogService)

	// 初始化中间件
	authMiddleware := middleware.NewAuthMiddleware(cfg)
	rateLimiter := middleware.NewRateLimiter(redisClient, cfg.Server.RateLimit)

	// 设置Gin模式
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API路由
	v1 := router.Group("/v1")
	v1.Use(rateLimiter.RateLimit())
	v1.Use(authMiddleware.JWTAuth())
	{
		// 订阅类型和状态
		v1.GET("/refs/subscription-types", refsHandler.GetSubscriptionTypes)
		v1.GET("/refs/subscription-statuses", refsHandler.GetSubscriptionStatuses)

		// 订阅管理
		v1.GET("/subscriptions", subscriptionHandler.GetSubscriptions)
		v1.POST("/subscriptions", subscriptionHandler.CreateSubscription)
		v1.GET("/subscriptions/:key", subscriptionHandler.GetSubscription)
		v1.GET("/subscriptions/:key/versions/:version", subscriptionHandler.GetSubscription)
		v1.PUT("/subscriptions/:key/versions/:version", subscriptionHandler.UpdateSubscription)
		v1.PATCH("/subscriptions/:key/versions/:version/status", subscriptionHandler.UpdateSubscriptionStatus)
		v1.DELETE("/subscriptions/:key/versions/:version", subscriptionHandler.DeleteSubscription)

		// 订阅执行
		v1.POST("/subscriptions/:key/execute", subscriptionHandler.ExecuteSubscription)
		v1.POST("/subscriptions/:key/versions/:version/execute", subscriptionHandler.ExecuteSubscription)

		// 统计查询
		v1.GET("/subscriptions/stats", subscriptionHandler.GetStats)

		// 操作日志
		v1.GET("/operation-logs", operationLogHandler.GetOperationLogs)
	}

	// 加载HTML模板
	router.LoadHTMLGlob("web/templates/*")

	// Web UI路由（简单的基础认证）
	webUI := router.Group("/admin")
	webUI.Use(authMiddleware.BasicAuth())
	{
		webUI.Static("/static", "./web/static")

		webUI.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "BI Subscription Management",
			})
		})

		webUI.GET("/subscriptions", func(c *gin.Context) {
			c.HTML(http.StatusOK, "subscriptions.html", gin.H{
				"title": "Subscription Management",
			})
		})

		webUI.GET("/stats", func(c *gin.Context) {
			c.HTML(http.StatusOK, "stats.html", gin.H{
				"title": "Statistics",
			})
		})

		webUI.GET("/operation-logs", func(c *gin.Context) {
			c.HTML(http.StatusOK, "operation_logs.html", gin.H{
				"title": "Operation Logs",
			})
		})
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
	}

	// 启动服务器
	go func() {
		log.Printf("Server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func initDB(config *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return db, nil
}
