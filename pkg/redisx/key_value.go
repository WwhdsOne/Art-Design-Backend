package redisx

import (
	"context"
	"go.uber.org/zap"
	"time"
)

// Set 方法用于设置 Redis 键值对
func (r *RedisWrapper) Set(id, value string, duration time.Duration) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.operationTimeout)
	defer cancelFunc()
	err = r.client.Set(timeout, id, value, duration).Err()
	if err != nil {
		return
	}
	return
}

// Get 方法用于获取 Redis 键对应的值
func (r *RedisWrapper) Get(key string) (val string, err error) {

	ctx, cancelFunc := context.WithTimeout(context.Background(), r.operationTimeout)
	defer cancelFunc()

	val, err = r.client.Get(ctx, key).Result()
	r.statsChan <- statsRecord{Key: key, IsHit: err == nil}
	if err != nil {
		return
	}
	return
}

// Del 方法用于删除 Redis 键
func (r *RedisWrapper) Del(key string) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.operationTimeout)
	defer cancelFunc()
	err = r.client.Del(timeout, key).Err()
	if err != nil {
		return
	}
	return
}

// PipelineSet 方法用于批量设置 Redis 键值对
func (r *RedisWrapper) PipelineSet(keyVal [][2]string, duration time.Duration) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.operationTimeout)
	defer cancelFunc()
	tx := r.client.Pipeline()
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
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.operationTimeout)
	defer cancelFunc()
	tx := r.client.Pipeline()
	for _, kv := range key {
		tx.Del(timeout, kv)
	}
	_, err = tx.Exec(timeout)
	if err != nil {
		return
	}
	return
}
