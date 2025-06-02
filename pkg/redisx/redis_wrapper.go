package redisx

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

// RedisWrapper 结构体用于封装 Redis 客户端和默认操作超时时间
type RedisWrapper struct {
	client           *redis.Client
	operationTimeout time.Duration // 操作超时时间

	hitCountMap   map[string]*atomic.Uint64
	totalCountMap map[string]*atomic.Uint64
	countLock     sync.RWMutex

	statsChan chan statsRecord

	scriptShaMap map[string]string // map[string] => string（脚本内容 -> SHA1）
	scriptLock   sync.RWMutex
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
	rw.hitCountMap = make(map[string]*atomic.Uint64)
	rw.totalCountMap = make(map[string]*atomic.Uint64)
	rw.scriptShaMap = make(map[string]string)
	// 读取持久化统计数据
	err := rw.LoadStatsFromRedis()
	if err != nil {
		zap.L().Fatal("无法读取redis键统计数据", zap.Error(err))
	}
	// 启动统计处理器
	go rw.statProcessor()
	// 启动命中率日志
	go rw.StartHitRateLogger(hitRateLogInterval)
	// 启动统计持久化
	go rw.SaveStatsToRedis(saveStatsInterval)
	return rw
}
