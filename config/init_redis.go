package config

import (
	"Art-Design-Backend/pkg/redisx"
	"Art-Design-Backend/pkg/utils"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Redis struct {
	Host               string `yaml:"host"`                  // 地址
	Port               string `yaml:"port"`                  // 端口
	Password           string `yaml:"password"`              // 密码（如果没有密码则为空）
	DB                 int    `yaml:"db"`                    // 数据库编号
	OperationTimeout   string `yaml:"operation-timeout"`     // 操作超时时间
	HitRateLogInterval string `yaml:"hit-rate-log-interval"` // 命中率日志间隔
	SaveStatsInterval  string `yaml:"save-stats-interval"`   // 保存统计信息间隔
}

func NewRedis(cfg *Config) *redisx.RedisWrapper {
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
