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

// Scan 方法用于扫描 Redis 键
// prefix 为键的前缀，cursor 为游标，count 为每次扫描的键数量
func (r *RedisWrapper) Scan(prefix string, cursor uint64, count int64) ([]string, uint64, error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.operationTimeout)
	defer cancelFunc()
	return r.client.Scan(timeout, cursor, prefix+"*", count).Result()
}

// DeleteByPrefix 方法用于根据前缀删除 Redis 键
// prefix 为键的前缀，count 为每次删除的键数量
func (r *RedisWrapper) DeleteByPrefix(prefix string, count int64) (err error) {
	var cursor uint64
	var keys []string

	for {
		// 扫描 keys
		keys, cursor, err = r.Scan(prefix, cursor, count)
		if err != nil {
			return
		}

		// 每批次创建并释放独立的 context
		if len(keys) > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), r.operationTimeout)
			err = r.client.Del(ctx, keys...).Err() // 忽略失败
			if err != nil {
				zap.L().Error("删除缓存失败", zap.Error(err))
			}
			cancel() // ✅ 必须释放
		}

		if cursor == 0 {
			break
		}
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
