package middleware

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/redisx"
	"gorm.io/gorm"
)

type Middlewares struct {
	Db     *gorm.DB             // 数据库
	Redis  *redisx.RedisWrapper // redis
	Jwt    *jwt.JWT             // jwt
	Config *config.Middleware   // 配置
}

func NewMiddlewares(
	db *gorm.DB,
	redis *redisx.RedisWrapper,
	jwt *jwt.JWT,
	c *config.Middleware,
) *Middlewares {
	return &Middlewares{
		Db:     db,
		Redis:  redis,
		Jwt:    jwt,
		Config: c,
	}
}
