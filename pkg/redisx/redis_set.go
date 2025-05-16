package redisx

import (
	"context"
	"go.uber.org/zap"
)

func (r *RedisWrapper) SMembers(key string) (val []string) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.OperationTimeout)
	defer cancelFunc()
	val, err := r.Client.SMembers(timeout, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			// 这里省略 return，默认返回 val 的零值 ""
		} else {
			zap.L().Error(err.Error())
		}
	}
	return
}

func (r *RedisWrapper) SAdd(key string, vals ...interface{}) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.OperationTimeout)
	defer cancelFunc()
	_, err = r.Client.SAdd(timeout, key, vals...).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	return
}
