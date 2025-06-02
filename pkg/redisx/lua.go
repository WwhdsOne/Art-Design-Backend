package redisx

import (
	"context"
	"strings"
)

func isNoScriptErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "NOSCRIPT")
}

// Eval 执行 Lua 脚本（自动缓存 SHA，自动降级）
func (r *RedisWrapper) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.operationTimeout)
	defer cancel()

	r.scriptLock.RLock()
	sha, ok := r.scriptShaMap[script]
	r.scriptLock.RUnlock()

	// 优先尝试使用缓存的 SHA
	if ok {
		cmd := r.client.EvalSha(ctx, sha, keys, args...)
		if err := cmd.Err(); err == nil {
			return cmd.Val(), nil
		} else if !isNoScriptErr(err) {
			return nil, err
		}
		// 否则 fallback 到 Eval
	}

	// 直接执行脚本（Eval）
	cmd := r.client.Eval(ctx, script, keys, args...)
	if err := cmd.Err(); err != nil {
		return nil, err
	}

	// 尝试加载 SHA 并缓存
	go func() {
		// 后台并行加载脚本 SHA（非阻塞主流程）
		ctx, cancel := context.WithTimeout(context.Background(), r.operationTimeout)
		defer cancel()

		if sha, err := r.client.ScriptLoad(ctx, script).Result(); err == nil {
			r.scriptLock.Lock()
			r.scriptShaMap[script] = sha
			r.scriptLock.Unlock()
		}
	}()

	return cmd.Val(), nil
}
