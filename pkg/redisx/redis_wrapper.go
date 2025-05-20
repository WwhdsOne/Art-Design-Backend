package redisx

import (
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

// RedisWrapper 结构体用于封装 Redis 客户端和默认操作超时时间
type RedisWrapper struct {
	client           *redis.Client
	operationTimeout time.Duration // 操作超时时间

	hitCountMap   sync.Map
	totalCountMap sync.Map

	statChan chan statRecord
}

type statRecord struct {
	Key   string
	IsHit bool
}

func NewRedisWrapper(client *redis.Client, timeout time.Duration, hitRateLogInterval time.Duration) *RedisWrapper {
	rw := &RedisWrapper{
		client:           client,
		operationTimeout: timeout,
		hitCountMap:      sync.Map{},
		totalCountMap:    sync.Map{},
		statChan:         make(chan statRecord, 1000), // 有缓冲，避免阻塞
	}
	// 启动统计处理器
	go rw.statProcessor()
	// 启动命中率日志
	go rw.StartHitRateLogger(hitRateLogInterval)
	return rw
}
