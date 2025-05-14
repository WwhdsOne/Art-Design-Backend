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
	Addr             string `yaml:"addr"`              // 地址
	Port             string `yaml:"port"`              // 端口
	Password         string `yaml:"password"`          // 密码（如果没有密码则为空）
	DB               int    `yaml:"db"`                // 数据库编号
	PreKey           string `yaml:"preKey"`            // 前缀
	OperationTimeout string `yaml:"operation-timeout"` // 操作超时时间
}

func NewRedis(cfg *Config) *redisx.RedisWrapper {
	r := cfg.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Addr, r.Port),
		Password: r.Password,
		DB:       r.DB,
	})
	// 检查连接是否成功
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		zap.L().Fatal("Redis 连接失败")
	}
	return &redisx.RedisWrapper{
		Client:           client,
		OperationTimeout: utils.ParseDuration(r.OperationTimeout),
	}
}
