package initialize

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/global"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) *redis.Client {
	r := cfg.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Addr, r.Port),
		Password: r.Password,
		DB:       r.DB,
	})

	// 检查连接是否成功
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.Logger.Fatal("Failed to connect to Redis")
	}
	// 自动迁移
	return client
}
