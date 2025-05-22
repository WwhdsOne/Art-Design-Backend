package redisx

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// StartHitRateLogger 启动一个协程定时打印命中率
func (r *RedisWrapper) StartHitRateLogger(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println()
		fmt.Println("┌────────────────────────────┬────────────┬────────────┬────────────┐")
		fmt.Println("│ Redis Key Prefix           │ Hit Count  │ Total Req  │ Hit Rate % │")
		fmt.Println("├────────────────────────────┼────────────┼────────────┼────────────┤")

		r.totalCountMap.Range(func(key, totalVal any) bool {
			category := key.(string)
			total := atomic.LoadUint64(totalVal.(*uint64))

			hitVal, ok := r.hitCountMap.Load(category)
			var hit uint64
			if ok {
				hit = atomic.LoadUint64(hitVal.(*uint64))
			}

			hitRate := 0.0
			if total > 0 {
				hitRate = float64(hit) / float64(total) * 100
			}

			fmt.Printf("│ %-26s │ %-10d │ %-10d │ %9.2f%% │\n", category, hit, total, hitRate)
			return true
		})

		fmt.Println("└────────────────────────────┴────────────┴────────────┴────────────┘")
	}
}

// 统计处理
func (r *RedisWrapper) statProcessor() {
	for stat := range r.statChan {
		category := getKeyCategory(stat.Key)
		// 增加总请求数
		r.incrMapCounter(&r.totalCountMap, category)

		// 命中时增加命中数
		if stat.IsHit {
			r.incrMapCounter(&r.hitCountMap, category)
		}
	}
}

func getKeyCategory(key string) string {
	idx := strings.LastIndex(key, ":")
	if idx == -1 {
		return key
	}
	return key[:idx+1]
}

// 原子计数
func (r *RedisWrapper) incrMapCounter(m *sync.Map, category string) {
	counterPtr, _ := m.LoadOrStore(category, new(uint64))
	atomic.AddUint64(counterPtr.(*uint64), 1)
}
