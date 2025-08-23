package middleware

import (
	"Art-Design-Backend/pkg/constant/rediskey"
	"Art-Design-Backend/pkg/result"
	"embed"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//go:embed redis-ratelimit.lua
var luaScript embed.FS

// RedisRateLimitMiddleware Redis 限流器
// windowSec: 窗口大小,单位秒
// maxReq: 最大请求数
func (m *Middlewares) RedisRateLimitMiddleware(windowSec int8, maxReq int8) gin.HandlerFunc {
	lua, err := luaScript.ReadFile("redis-ratelimit.lua")
	if err != nil {
		zap.L().Fatal("Redis 限流器Lua脚本读取失败")
	}
	script := string(lua)
	return func(c *gin.Context) {
		key := fmt.Sprintf(rediskey.RateLimiter+"%s", c.ClientIP())
		now := time.Now().Unix()

		res, err := m.Redis.Eval(script, []string{key}, windowSec, maxReq, now)

		if err != nil {
			// Redis 出错时可选：放行或拦截
			result.FailWithMessage("Redis 限流器出错", c)
			c.Abort()
			return
		}

		allowed, _ := res.(int64)
		if allowed <= 0 {
			result.FailWithMessage("请求过于频繁", c)
			c.Abort()
			return
		}

		c.Next()
	}
}
