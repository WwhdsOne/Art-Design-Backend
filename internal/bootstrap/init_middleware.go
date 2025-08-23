package bootstrap

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/redisx"

	"gorm.io/gorm"
)

func InitMiddleware(
	cfg *config.Config,
	db *gorm.DB,
	redis *redisx.RedisWrapper,
	jwt *jwt.JWT,
) *middleware.Middlewares {
	return &middleware.Middlewares{
		Config: &cfg.Middleware,
		Db:     db,
		Redis:  redis,
		Jwt:    jwt,
	}
}
