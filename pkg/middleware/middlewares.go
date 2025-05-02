package middleware

import (
	"Art-Design-Backend/pkg/jwt"
	"Art-Design-Backend/pkg/redisx"
	"gorm.io/gorm"
)

type Middlewares struct {
	Db    *gorm.DB             // 数据库
	Redis *redisx.RedisWrapper // redis
	Jwt   *jwt.JWT             // jwt
}
