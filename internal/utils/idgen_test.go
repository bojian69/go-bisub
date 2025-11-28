package utils

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitSnowflake(t *testing.T) {
	tests := []struct {
		name    string
		nodeID  int64
		wantErr bool
	}{
		{
			name:    "valid_node_id",
			nodeID:  1,
			wantErr: false,
		},
		{
			name:    "max_node_id",
			nodeID:  1023,
			wantErr: false,
		},
		{
			name:    "invalid_node_id_negative",
			nodeID:  -1,
			wantErr: true,
		},
		{
			name:    "invalid_node_id_too_large",
			nodeID:  1024,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置全局变量
			snowflakeNode = nil
			once = sync.Once{}
			
			err := InitSnowflake(tt.nodeID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, snowflakeNode)
			}
		})
	}
}

func TestGenerateID(t *testing.T) {
	// 初始化
	snowflakeNode = nil
	once = sync.Once{}
	err := InitSnowflake(1)
	require.NoError(t, err)

	// 测试ID生成
	id := GenerateID()
	assert.Greater(t, id, int64(0))

	// 测试ID唯一性
	ids := make(map[int64]bool)
	for i := 0; i < 10000; i++ {
		id := GenerateID()
		assert.False(t, ids[id], "duplicate ID generated: %d", id)
		ids[id] = true
	}
}

func TestGenerateIDString(t *testing.T) {
	snowflakeNode = nil
	once = sync.Once{}
	err := InitSnowflake(1)
	require.NoError(t, err)

	idStr := GenerateIDString()
	assert.NotEmpty(t, idStr)
	assert.Regexp(t, `^\d+$`, idStr)
}

func TestGenerateUUIDv7(t *testing.T) {
	uuid1 := GenerateUUIDv7()
	uuid2 := GenerateUUIDv7()
	
	assert.NotEmpty(t, uuid1)
	assert.NotEmpty(t, uuid2)
	assert.NotEqual(t, uuid1, uuid2)
	
	// UUID v7格式验证
	assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, uuid1)
}

func TestGenerateUUIDv4(t *testing.T) {
	uuid1 := GenerateUUIDv4()
	uuid2 := GenerateUUIDv4()
	
	assert.NotEmpty(t, uuid1)
	assert.NotEmpty(t, uuid2)
	assert.NotEqual(t, uuid1, uuid2)
	
	// UUID v4格式验证
	assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, uuid1)
}

func TestGenerateIDs(t *testing.T) {
	snowflakeNode = nil
	once = sync.Once{}
	err := InitSnowflake(1)
	require.NoError(t, err)

	count := 100
	ids := GenerateIDs(count)
	
	assert.Len(t, ids, count)
	
	// 检查唯一性
	idMap := make(map[int64]bool)
	for _, id := range ids {
		assert.False(t, idMap[id], "duplicate ID in batch: %d", id)
		idMap[id] = true
	}
}

func TestConcurrentIDGeneration(t *testing.T) {
	snowflakeNode = nil
	once = sync.Once{}
	err := InitSnowflake(1)
	require.NoError(t, err)

	var wg sync.WaitGroup
	ids := sync.Map{}
	duplicates := int32(0)
	
	// 并发生成ID
	goroutines := 100
	idsPerGoroutine := 1000
	
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id := GenerateID()
				if _, exists := ids.LoadOrStore(id, true); exists {
					atomic.AddInt32(&duplicates, 1)
				}
			}
		}()
	}
	
	wg.Wait()
	assert.Equal(t, int32(0), duplicates, "found duplicate IDs in concurrent generation")
}

func TestIDPool(t *testing.T) {
	snowflakeNode = nil
	once = sync.Once{}
	err := InitSnowflake(1)
	require.NoError(t, err)

	poolSize := 10
	pool := NewIDPool(poolSize)
	
	// 测试从池中获取ID
	ids := make(map[int64]bool)
	for i := 0; i < poolSize*2; i++ {
		id := pool.Get()
		assert.Greater(t, id, int64(0))
		assert.False(t, ids[id], "duplicate ID from pool: %d", id)
		ids[id] = true
	}
}

func TestParseSnowflakeID(t *testing.T) {
	// 保存原始状态
	origNode := snowflakeNode
	origOnce := once
	
	// 测试结束后恢复
	t.Cleanup(func() {
		snowflakeNode = origNode
		once = origOnce
	})
	
	// 重置并初始化
	snowflakeNode = nil
	once = sync.Once{}
	err := InitSnowflake(1)
	require.NoError(t, err)

	id := GenerateID()
	timestamp := ParseSnowflakeID(id)
	
	// 时间戳应该接近当前时间
	now := time.Now()
	diff := now.Sub(timestamp)
	assert.Less(t, diff, time.Second, "timestamp should be close to current time")
}

func TestValidateSnowflakeID(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "valid_id",
			id:      1234567890123456789,
			wantErr: false,
		},
		{
			name:    "zero_id",
			id:      0,
			wantErr: true,
		},
		{
			name:    "negative_id",
			id:      -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSnowflakeID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateIDWithFallback(t *testing.T) {
	// 测试未初始化时的降级
	snowflakeNode = nil
	once = sync.Once{}
	
	id := GenerateIDWithFallback()
	assert.NotEmpty(t, id)
	
	// 应该是UUID格式
	assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, id)
}

func BenchmarkGenerateID(b *testing.B) {
	snowflakeNode = nil
	once = sync.Once{}
	err := InitSnowflake(1)
	require.NoError(b, err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GenerateID()
		}
	})
}

func BenchmarkGenerateUUIDv7(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GenerateUUIDv7()
		}
	})
}

func BenchmarkIDPool(b *testing.B) {
	snowflakeNode = nil
	once = sync.Once{}
	err := InitSnowflake(1)
	require.NoError(b, err)

	pool := NewIDPool(1000)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.Get()
		}
	})
}