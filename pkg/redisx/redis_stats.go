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

		r.countLock.RLock()
		for k, v := range r.hitCountMap {
			data.HitCounts[k] = v.Load()
		}
		for k, v := range r.totalCountMap {
			data.TotalCounts[k] = v.Load()
		}
		r.countLock.RUnlock()

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

	r.countLock.Lock()
	for k, v := range data.HitCounts {
		counter := new(atomic.Uint64)
		counter.Store(v)
		r.hitCountMap[k] = counter
	}
	for k, v := range data.TotalCounts {
		counter := new(atomic.Uint64)
		counter.Store(v)
		r.totalCountMap[k] = counter
	}
	r.countLock.Unlock()
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

		r.countLock.RLock()
		for category, totalCounter := range r.totalCountMap {
			total := totalCounter.Load()
			hit := uint64(0)
			if hitCounter, ok := r.hitCountMap[category]; ok {
				hit = hitCounter.Load()
			}
			hitRate := 0.0
			if total > 0 {
				hitRate = float64(hit) / float64(total) * 100
			}
			fmt.Printf("│ %-26s │ %-10d │ %-10d │ %9.2f%% │\n", category, hit, total, hitRate)
		}
		r.countLock.RUnlock()

		fmt.Println("└────────────────────────────┴────────────┴────────────┴────────────┘")
	}
}

// 统计处理
func (r *RedisWrapper) statProcessor() {
	for stat := range r.statsChan {
		category := getKeyCategory(stat.Key)
		// 增加总请求数
		r.incrMapCounter(r.totalCountMap, category)

		// 命中时增加命中数
		if stat.IsHit {
			r.incrMapCounter(r.hitCountMap, category)
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
func (r *RedisWrapper) incrMapCounter(m map[string]*atomic.Uint64, category string) {
	// 尝试读锁获取已有 counter
	r.countLock.RLock()
	counter, ok := m[category]
	r.countLock.RUnlock()

	// 已有则直接加，无需写锁
	if ok {
		counter.Add(1)
		return
	}

	// 否则写入新 counter
	r.countLock.Lock()
	// 双重检查，避免并发写入
	counter, ok = m[category]
	if !ok {
		counter = new(atomic.Uint64)
		m[category] = counter
	}
	r.countLock.Unlock()

	counter.Add(1)
}
