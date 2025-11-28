package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	snowflakeNode *snowflake.Node
	once          sync.Once
)

// InitSnowflake 初始化Snowflake节点
func InitSnowflake(nodeID int64) error {
	var err error
	once.Do(func() {
		snowflakeNode, err = snowflake.NewNode(nodeID)
		if err != nil {
			logrus.WithError(err).WithField("node_id", nodeID).Error("failed to initialize snowflake node")
		} else {
			logrus.WithField("node_id", nodeID).Info("snowflake node initialized successfully")
		}
	})
	return err
}

// GenerateID 生成分布式ID (Snowflake)
func GenerateID() int64 {
	if snowflakeNode == nil {
		logrus.Fatal("snowflake node not initialized, call InitSnowflake first")
	}
	return snowflakeNode.Generate().Int64()
}

// GenerateIDString 生成字符串格式ID (Snowflake)
func GenerateIDString() string {
	if snowflakeNode == nil {
		logrus.Fatal("snowflake node not initialized, call InitSnowflake first")
	}
	return snowflakeNode.Generate().String()
}

// GenerateUUIDv7 生成UUID v7 (时间排序)
func GenerateUUIDv7() string {
	return uuid.Must(uuid.NewV7()).String()
}

// GenerateUUIDv4 生成UUID v4 (随机)
func GenerateUUIDv4() string {
	return uuid.New().String()
}

// GenerateIDs 批量生成ID
func GenerateIDs(count int) []int64 {
	if snowflakeNode == nil {
		logrus.Fatal("snowflake node not initialized")
	}
	
	ids := make([]int64, count)
	for i := 0; i < count; i++ {
		ids[i] = snowflakeNode.Generate().Int64()
	}
	return ids
}

// GenerateIDWithFallback 生成ID，失败时降级到UUID
func GenerateIDWithFallback() string {
	defer func() {
		if r := recover(); r != nil {
			logrus.WithField("panic", r).Error("snowflake generation failed, fallback to UUID")
		}
	}()
	
	if snowflakeNode != nil {
		return GenerateIDString()
	}
	
	logrus.Warn("snowflake not available, using UUID v7")
	return GenerateUUIDv7()
}

// IDPool ID预生成池
type IDPool struct {
	pool chan int64
	size int
	mu   sync.RWMutex
}

// NewIDPool 创建ID池
func NewIDPool(size int) *IDPool {
	p := &IDPool{
		pool: make(chan int64, size),
		size: size,
	}
	p.fill()
	return p
}

// Get 从池中获取ID
func (p *IDPool) Get() int64 {
	select {
	case id := <-p.pool:
		go p.refill() // 异步补充
		return id
	default:
		// 池空时直接生成
		logrus.Debug("ID pool empty, generating new ID")
		return GenerateID()
	}
}

// fill 填充ID池
func (p *IDPool) fill() {
	for i := 0; i < p.size; i++ {
		select {
		case p.pool <- GenerateID():
		default:
			return // 池已满
		}
	}
}

// refill 异步补充ID池
func (p *IDPool) refill() {
	p.mu.Lock()
	defer p.mu.Unlock()
	
	// 补充到一半容量
	target := p.size / 2
	current := len(p.pool)
	
	for i := current; i < target; i++ {
		select {
		case p.pool <- GenerateID():
		default:
			return
		}
	}
}

// ParseSnowflakeID 解析Snowflake ID获取时间戳
func ParseSnowflakeID(id int64) time.Time {
	// Snowflake ID结构: 1位符号位 + 41位时间戳 + 10位机器ID + 12位序列号
	// 时间戳是相对于2006-03-21 20:50:14 UTC的毫秒数
	timestamp := (id >> 22) + 1288834974657 // snowflake epoch
	return time.Unix(timestamp/1000, (timestamp%1000)*1000000)
}

// ValidateSnowflakeID 验证Snowflake ID格式
func ValidateSnowflakeID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid snowflake ID: %d", id)
	}
	
	// 检查时间戳是否合理 (不能是未来时间)
	timestamp := ParseSnowflakeID(id)
	if timestamp.After(time.Now().Add(time.Minute)) {
		return fmt.Errorf("snowflake ID timestamp is in the future: %v", timestamp)
	}
	
	return nil
}