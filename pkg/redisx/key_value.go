package redisx

import (
	"context"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

// Set 方法用于设置 Redis 键值对
func (r *RedisWrapper) Set(id, value string, duration time.Duration) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.OperationTimeout)
	defer cancelFunc()
	err = r.Client.Set(timeout, id, value, duration).Err()
	if err != nil {
		return
	}
	return
}

// Get 方法用于获取 Redis 键对应的值
func (r *RedisWrapper) Get(key string) (val string, err error) {
	atomic.AddUint64(&r.totalCount, 1)

	ctx, cancelFunc := context.WithTimeout(context.Background(), r.OperationTimeout)
	defer cancelFunc()

	val, err = r.Client.Get(ctx, key).Result()
	if err != nil {
		return
	}

	// 命中缓存
	atomic.AddUint64(&r.hitCount, 1)
	return
}

// Del 方法用于删除 Redis 键
func (r *RedisWrapper) Del(key string) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.OperationTimeout)
	defer cancelFunc()
	err = r.Client.Del(timeout, key).Err()
	if err != nil {
		return
	}
	return
}

// PipelineSet 方法用于批量设置 Redis 键值对
func (r *RedisWrapper) PipelineSet(keyVal [][2]string, duration time.Duration) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.OperationTimeout)
	defer cancelFunc()
	tx := r.Client.Pipeline()
	for _, kv := range keyVal {
		tx.Set(timeout, kv[0], kv[1], duration)
	}
	_, err = tx.Exec(timeout)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return
}

// PipelineDel 方法用于批量删除 Redis 键
func (r *RedisWrapper) PipelineDel(key []string) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.OperationTimeout)
	defer cancelFunc()
	tx := r.Client.Pipeline()
	for _, kv := range key {
		tx.Del(timeout, kv)
	}
	_, err = tx.Exec(timeout)
	if err != nil {
		return
	}
	return
}
