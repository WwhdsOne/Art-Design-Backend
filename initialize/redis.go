package initialize

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/global"
	"Art-Design-Backend/pkg/redisx"
	"Art-Design-Backend/pkg/utils"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) redisx.RedisWrapper {
	r := cfg.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Addr, r.Port),
		Password: r.Password,
		DB:       r.DB,
	})
	// 检查连接是否成功
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.Logger.Fatal("Redis 连接失败")
	}
	// 自动迁移
	return redisx.RedisWrapper{
		Client:         client,
		DefaultTimeout: utils.ParseDuration(r.OperationTimeout),
	}
}
