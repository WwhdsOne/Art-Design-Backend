package controller

import (
	"Art-Design-Backend/internal/service"
	"github.com/google/wire"
)

var ControllersProvider = wire.NewSet(
	AuthCtrlProvider,
	UserCtrlProvider,
	MenuCtrlProvider,
	RoleCtrlProvider,
	DigitPredictProvider,
)

var AuthCtrlProvider = wire.NewSet(
	NewAuthController,
	wire.Struct(new(service.AuthService), "*"),
)

var UserCtrlProvider = wire.NewSet(
	NewUserController,
	wire.Struct(new(service.UserService), "*"),
)

var MenuCtrlProvider = wire.NewSet(
	NewMenuController,
	wire.Struct(new(service.MenuService), "*"),
)

var RoleCtrlProvider = wire.NewSet(
	NewRoleController,
	wire.Struct(new(service.RoleService), "*"),
)

var DigitPredictProvider = wire.NewSet(
	NewDigitPredictController,
	wire.Struct(new(service.DigitPredictService), "*"),
)
