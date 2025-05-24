-- redis-ratelimit.lua
-- KEYS[1]: key（比如 IP 或用户ID）
-- ARGV[1]: 窗口大小（秒）
-- ARGV[2]: 最大请求数
-- ARGV[3]: 当前时间戳（秒）

local key = KEYS[1]
local window = tonumber(ARGV[1])
local limit = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

-- 清理旧的时间戳
redis.call("ZREMRANGEBYSCORE", key, 0, now - window)

-- 获取剩余请求数
local count = redis.call("ZCARD", key)
if count >= limit then
    return 0
end

-- 添加当前请求时间戳
redis.call("ZADD", key, now, now)
redis.call("EXPIRE", key, window + 1)

return 1
