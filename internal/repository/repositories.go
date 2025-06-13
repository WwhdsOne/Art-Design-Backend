package repository

import (
	"Art-Design-Backend/internal/repository/cache"
	"Art-Design-Backend/internal/repository/db"
	"github.com/google/wire"
)

var RepositoriesProvider = wire.NewSet(
	db.NewUserDB,
	db.NewMenuRepository,
	db.NewRoleRepository,
	db.NewGormTransactionManager,
	db.NewRoleMenusRepository,
	db.NewUserRolesRepository,
	db.NewDigitPredictRepository,
	db.NewAIModelRepository,
	cache.NewAIModelCache,
	cache.NewAuthCache,
)
