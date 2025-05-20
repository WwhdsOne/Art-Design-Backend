package redisx

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// StartHitRateLogger 启动一个协程定时打印命中率
func (r *RedisWrapper) StartHitRateLogger(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
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
	}()
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

// 使用缓冲池，减少内存分配
// 指定最长长度，避免频繁分配内存
// 128字节兼顾性能和空间
var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 128) // 128字节
	},
}

// 根据 key 获取分类
func getKeyCategory(key string) string {
	b := bufferPool.Get().([]byte)
	b = append(b[:0], key...) // 重置并复用切片

	idx := bytes.LastIndexByte(b, ':')
	if idx == -1 {
		bufferPool.Put(b)
		return key
	}

	result := string(b[:idx+1])
	bufferPool.Put(b)
	return result
}

// 原子计数
func (r *RedisWrapper) incrMapCounter(m *sync.Map, category string) {
	counterPtr, _ := m.LoadOrStore(category, new(uint64))
	atomic.AddUint64(counterPtr.(*uint64), 1)
}
