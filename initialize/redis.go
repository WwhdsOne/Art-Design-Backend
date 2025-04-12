package initialize

import (
	"Art-Design-Backend/config"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) *redis.Client {
	r := cfg.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Addr, r.Port),
		Password: r.Password, // no password set
		DB:       r.DB,       // use default DB
	})
	// 自动迁移
	return client
}
