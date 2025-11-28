package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	redis *redis.Client
	limit int
}

func NewRateLimiter(redisClient *redis.Client, limit int) *RateLimiter {
	return &RateLimiter{
		redis: redisClient,
		limit: limit,
	}
}

// RateLimit 限流中间件
func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate_limit:%s", c.ClientIP())

		ctx := context.Background()

		// 使用滑动窗口算法
		now := time.Now().Unix()
		window := int64(60) // 1分钟窗口

		pipe := rl.redis.Pipeline()

		// 移除过期的请求
		pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(now-window, 10))

		// 添加当前请求
		pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})

		// 获取当前窗口内的请求数
		pipe.ZCard(ctx, key)

		// 设置过期时间
		pipe.Expire(ctx, key, time.Duration(window)*time.Second)

		results, err := pipe.Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Rate limit check failed",
			})
			c.Abort()
			return
		}

		// 获取请求数
		count := results[2].(*redis.IntCmd).Val()

		// 设置响应头
		c.Header("X-RateLimit-Limit", strconv.Itoa(rl.limit))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(int64(rl.limit)-count, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(now+window, 10))

		if count > int64(rl.limit) {
			c.Header("Retry-After", "60")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "RATE_LIMITED",
				"message": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
