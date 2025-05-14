package repository

import "github.com/google/wire"

var RepositoriesProvider = wire.NewSet(
	NewUserRepository,
	NewMenuRepository,
	NewRoleRepository,
	NewGormTransactionManager,
)
