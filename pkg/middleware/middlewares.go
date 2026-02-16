package middleware

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/internal/repository/db"
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/redisx"
)

type Middlewares struct {
	Redis          *redisx.RedisWrapper // redis
	Jwt            *jwt.JWT             // jwt
	Config         *config.Middleware   // 配置
	OperationLogDB *db.OperationLogDB   // 操作日志
}
