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
	cache.NewAIProviderCache,
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
	db.NewAIProviderDB,
	db.NewAIAgentDB,
	db.NewAgentFileDB,
	db.NewFileChunkDB,
	db.NewChunkVectorDB,
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
	wire.Struct(new(AIProviderRepo), "*"),
	wire.Struct(new(AIAgentRepo), "*"),
)
