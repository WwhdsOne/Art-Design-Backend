package redisx

import (
	"Art-Design-Backend/pkg/errors"
	"context"

	"go.uber.org/zap"
)

func (r *RedisWrapper) SMembers(key string) (val []string) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.operationTimeout)
	defer cancelFunc()
	// 映射表默认永久存在
	val, err := r.client.SMembers(timeout, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			// 这里省略 return，默认返回 val 的零值 ""
		} else {
			zap.L().Error(err.Error())
		}
	}
	return
}

// SAdd 添加元素
// 永久存在
func (r *RedisWrapper) SAdd(key string, vals ...interface{}) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.operationTimeout)
	defer cancelFunc()
	_, err = r.client.SAdd(timeout, key, vals...).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	return
}

// SRem 删除元素
func (r *RedisWrapper) SRem(key string, members ...interface{}) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), r.operationTimeout)
	defer cancelFunc()
	_, err = r.client.SRem(timeout, key, members...).Result()
	if err != nil {
		return
	}
	return
}

// DelBySetMembers 根据集合 key，删除集合中每个成员对应的缓存键，最后删除集合自身
// 使用 Lua 脚本确保操作的原子性
func (r *RedisWrapper) DelBySetMembers(setKey string) (err error) {

	// Lua 脚本：获取集合成员，依次删除每个以成员为 key 的缓存，再删除集合自身
	script := `
        local members = redis.call('SMEMBERS', KEYS[1])
        if #members > 0 then
            for i, member in ipairs(members) do
                redis.call('DEL', member)
            end
        end
        return redis.call('DEL', KEYS[1])
    `

	// 执行 Lua 脚本
	_, err = r.Eval(script, []string{setKey})
	if err != nil {
		err = errors.WrapCacheError(err, "根据集合成员删除键失败")
		return
	}

	return
}
