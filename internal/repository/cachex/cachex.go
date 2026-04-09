package cachex

import (
	"Art-Design-Backend/pkg/redisx"
	"time"

	"github.com/bytedance/sonic"
)

// GetOrNil 从 Redis 反序列化单个对象，缓存未命中时返回 nil, redis.Nil
func GetOrNil[T any](r *redisx.RedisWrapper, key string) (*T, error) {
	val, err := r.Get(key)
	if err != nil {
		return nil, err
	}
	var result T
	if err = sonic.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetSlice 从 Redis 反序列化切片，缓存未命中时返回 nil, redis.Nil
func GetSlice[T any](r *redisx.RedisWrapper, key string) ([]T, error) {
	val, err := r.Get(key)
	if err != nil {
		return nil, err
	}
	var result []T
	if err = sonic.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Set 将对象序列化后写入 Redis
func Set[T any](r *redisx.RedisWrapper, key string, val T, ttl time.Duration) error {
	data, err := sonic.Marshal(val)
	if err != nil {
		return err
	}
	return r.Set(key, string(data), ttl)
}
