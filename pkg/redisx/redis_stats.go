package redisx

import (
	"Art-Design-Backend/pkg/constant/rediskey"
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type persistData struct {
	HitCounts   map[string]uint64 `json:"hit_counts"`
	TotalCounts map[string]uint64 `json:"total_counts"`
}

func (r *RedisWrapper) SaveStatsToRedis(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		data := persistData{
			HitCounts:   make(map[string]uint64),
			TotalCounts: make(map[string]uint64),
		}

		r.hitCountMap.Range(func(key, value any) bool {
			data.HitCounts[key.(string)] = atomic.LoadUint64(value.(*uint64))
			return true
		})
		r.totalCountMap.Range(func(key, value any) bool {
			data.TotalCounts[key.(string)] = atomic.LoadUint64(value.(*uint64))
			return true
		})

		jsonBytes, _ := sonic.Marshal(data)
		err := r.client.Set(context.Background(), rediskey.KeyStats, jsonBytes, 0).Err()
		if err != nil {
			zap.L().Error("Redis键统计数据存储失败", zap.Error(err))
		}
	}
}

func (r *RedisWrapper) LoadStatsFromRedis() error {
	val, err := r.client.Get(context.Background(), rediskey.KeyStats).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}

	var data persistData
	if err = sonic.Unmarshal([]byte(val), &data); err != nil {
		return err
	}

	for k, v := range data.HitCounts {
		count := new(uint64)
		atomic.StoreUint64(count, v)
		r.hitCountMap.Store(k, count)
	}
	for k, v := range data.TotalCounts {
		count := new(uint64)
		atomic.StoreUint64(count, v)
		r.totalCountMap.Store(k, count)
	}
	return nil
}

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
func (r *RedisWrapper) statsProcessor() {
	for stat := range r.statsChan {
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
	// 优先尝试 Load，避免不必要的开销
	if val, ok := m.Load(category); ok {
		atomic.AddUint64(val.(*uint64), 1)
		return
	}

	// 如果未找到，则尝试 LoadOrStore（双重检查）
	newCounter := new(uint64)
	actual, _ := m.LoadOrStore(category, newCounter)
	atomic.AddUint64(actual.(*uint64), 1)
}
