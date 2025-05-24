package redisx

import "context"

// todo后续可以优化为SHA加快速度

// Eval 执行lua脚本
func (r *RedisWrapper) Eval(script string, keys []string, args ...interface{}) (result interface{}, err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.operationTimeout)
	defer cancelFunc()
	result, err = r.client.Eval(timeout, script, keys, args...).Result()
	if err != nil {
		return
	}
	return
}
