package redisx

import (
	"context"
	"strings"
)

func isNoScriptErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "NOSCRIPT")
}

// Eval 执行 Lua 脚本（自动缓存 SHA，自动降级）
func (r *RedisWrapper) Eval(script string, keys []string, args ...interface{}) (result interface{}, err error) {
	timeout, cancel := context.WithTimeout(context.Background(), r.operationTimeout)
	defer cancel()

	// 尝试使用 EVALSHA
	if val, ok := r.scriptShaMap.Load(script); ok {
		sha := val.(string)
		cmd := r.client.EvalSha(timeout, sha, keys, args...)
		if cmd.Err() == nil {
			return cmd.Val(), nil
		}
		// 如果不是 NOSCRIPT 错误，直接返回错误
		if !isNoScriptErr(cmd.Err()) {
			return nil, cmd.Err()
		}
		// 否则降级到 Eval 继续执行
	}

	// 使用 Eval 执行脚本
	cmd := r.client.Eval(timeout, script, keys, args...)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	// 尝试加载脚本缓存 SHA
	if sha, shaErr := r.client.ScriptLoad(timeout, script).Result(); shaErr == nil {
		r.scriptShaMap.Store(script, sha)
	}

	return cmd.Val(), nil
}
