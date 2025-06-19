package repository

import (
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	"github.com/google/wire"
)

var RedisCacheProvider = wire.NewSet(
	cache.NewAuthCache,
	cache.NewRoleCache,
	cache.NewMenuCache,
	cache.NewUserCache,
	cache.NewAIModelCache,
)

var DBProvider = wire.NewSet(
	db.NewUserDB,
	db.NewMenuDB,
	db.NewRoleDB,
	db.NewGormTransactionManager,
	db.NewRoleMenusDB,
	db.NewUserRolesDB,
	db.NewDigitPredictDB,
	db.NewAIModelDB,
)

var RepositoriesProvider = wire.NewSet(
	DBProvider,
	RedisCacheProvider,
	NewUserRepo,
	NewAIModelRepo,
	NewRoleRepo,
	NewMenuRepo,
	NewAuthRepo,
	NewDigitPredictRepo,
)
