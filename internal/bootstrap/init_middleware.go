package bootstrap

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/internal/repository/db"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/redisx"
)

func InitMiddleware(
	cfg *config.Config,
	redis *redisx.RedisWrapper,
	jwtService *jwt.JWT,
	operationLogDB *db.OperationLogDB,
) *middleware.Middlewares {
	return &middleware.Middlewares{
		Config:         &cfg.Middleware,
		Redis:          redis,
		Jwt:            jwtService,
		OperationLogDB: operationLogDB,
	}
}
