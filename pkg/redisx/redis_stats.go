package redisx

import (
	"fmt"
	"sync/atomic"
	"time"
)

// HitRate 返回当前命中率
func (r *RedisWrapper) HitRate() float64 {
	total := atomic.LoadUint64(&r.totalCount)
	if total == 0 {
		return 0
	}
	hit := atomic.LoadUint64(&r.hitCount)
	return float64(hit) / float64(total)
}

// StartHitRateLogger 启动一个协程定时打印命中率
func (r *RedisWrapper) StartHitRateLogger(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			rate := r.HitRate()
			fmt.Printf("[RedisWrapper] Cache Hit Rate: %.2f%% (hits: %d / total: %d)\n",
				rate*100,
				atomic.LoadUint64(&r.hitCount),
				atomic.LoadUint64(&r.totalCount),
			)
		}
	}()
}
