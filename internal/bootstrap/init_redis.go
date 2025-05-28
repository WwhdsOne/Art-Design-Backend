package bootstrap

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/redisx"
	"Art-Design-Backend/pkg/utils"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func InitRedis(cfg *config.Config) *redisx.RedisWrapper {
	r := cfg.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: r.Password,
		DB:       r.DB,
	})
	// 检查连接是否成功
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		zap.L().Fatal("Redis 连接失败")
	}
	return redisx.NewRedisWrapper(client,
		utils.ParseDuration(r.OperationTimeout),
		utils.ParseDuration(r.HitRateLogInterval),
		utils.ParseDuration(r.SaveStatsInterval),
	)
}
