package redisx

import (
	"Art-Design-Backend/global"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var redisInstance RedisWrapper

// RedisWrapper 结构体用于封装 Redis 客户端和默认操作超时时间
type RedisWrapper struct {
	Client         *redis.Client // 假设 global.Redis 是 *redis.Client 类型，这里用 RedisClientType 替代
	DefaultTimeout time.Duration
}

func NewRedisWrapper(r RedisWrapper) {
	redisInstance = r
}

// Set 方法用于设置 Redis 键值对
func Set(id, value string, duration time.Duration) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), redisInstance.DefaultTimeout)
	defer cancelFunc()
	err = redisInstance.Client.Set(timeout, id, value, duration).Err()
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return
}

// Get 方法用于获取 Redis 键对应的值
func Get(key string) (val string) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), redisInstance.DefaultTimeout)
	defer cancelFunc()
	val, err := redisInstance.Client.Get(timeout, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			// 这里省略 return，默认返回 val 的零值 ""
		} else {
			global.Logger.Error(err.Error())
		}
	}
	return
}

// Delete 方法用于删除 Redis 键
func Delete(key string) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), redisInstance.DefaultTimeout)
	defer cancelFunc()
	err = redisInstance.Client.Del(timeout, key).Err()
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return
}

// PipelineSet 方法用于批量设置 Redis 键值对
func PipelineSet(keyVal [][2]string, duration time.Duration) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), redisInstance.DefaultTimeout)
	defer cancelFunc()
	tx := redisInstance.Client.Pipeline()
	for _, kv := range keyVal {
		tx.Set(timeout, kv[0], kv[1], duration)
	}
	_, err = tx.Exec(timeout)
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return
}

// PipelineDelete 方法用于批量删除 Redis 键
func PipelineDelete(key []string) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), redisInstance.DefaultTimeout)
	defer cancelFunc()
	tx := redisInstance.Client.Pipeline()
	for _, kv := range key {
		tx.Del(timeout, kv)
	}
	_, err = tx.Exec(timeout)
	if err != nil {
		global.Logger.Error(err.Error())
	}
	return
}
