package repository

import (
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	"github.com/google/wire"
)

var RedisCacheSet = wire.NewSet(
	cache.NewAuthCache,
	cache.NewRoleCache,
	cache.NewMenuCache,
	cache.NewUserCache,
	cache.NewAIModelCache,
)

var DBSet = wire.NewSet(
	db.NewUserDB,
	db.NewMenuDB,
	db.NewRoleDB,
	db.NewGormTransactionManager,
	db.NewRoleMenusDB,
	db.NewUserRolesDB,
	db.NewDigitPredictDB,
	db.NewAIModelDB,
)

var RepositorySet = wire.NewSet(
	DBSet,
	RedisCacheSet,
	wire.Struct(new(UserRepo), "*"),
	wire.Struct(new(AIModelRepo), "*"),
	wire.Struct(new(RoleRepo), "*"),
	wire.Struct(new(MenuRepo), "*"),
	wire.Struct(new(AuthRepo), "*"),
	wire.Struct(new(DigitPredictRepo), "*"),
)
