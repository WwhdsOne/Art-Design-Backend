package redisx

import (
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisWrapper 结构体用于封装 Redis 客户端和默认操作超时时间
type RedisWrapper struct {
	client           *redis.Client
	operationTimeout time.Duration // 操作超时时间

	hitCountMap   sync.Map
	totalCountMap sync.Map

	statsChan chan statsRecord

	scriptShaMap sync.Map // map[string] => string（脚本内容 -> SHA1）
}

type statsRecord struct {
	Key   string
	IsHit bool
}

func NewRedisWrapper(client *redis.Client, timeout time.Duration, hitRateLogInterval time.Duration, saveStatsInterval time.Duration) *RedisWrapper {
	rw := &RedisWrapper{
		client:           client,
		operationTimeout: timeout,
		statsChan:        make(chan statsRecord, 1000), // 有缓冲，避免阻塞
	}
	// 读取持久化统计数据
	err := rw.LoadStatsFromRedis()
	if err != nil {
		zap.L().Fatal("无法读取redis键统计数据", zap.Error(err))
	}
	// 启动统计处理器
	go rw.statsProcessor()
	// 启动命中率日志
	go rw.StartHitRateLogger(hitRateLogInterval)
	// 启动统计持久化
	go rw.SaveStatsToRedis(saveStatsInterval)
	return rw
}
