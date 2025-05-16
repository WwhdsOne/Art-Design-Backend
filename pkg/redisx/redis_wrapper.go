package redisx

import (
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisWrapper 结构体用于封装 Redis 客户端和默认操作超时时间
type RedisWrapper struct {
	Client           *redis.Client
	OperationTimeout time.Duration // 默认操作超时时间
}
