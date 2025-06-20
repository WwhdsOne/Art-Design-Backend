package controller

import (
	"Art-Design-Backend/internal/service"
	"github.com/google/wire"
)

var ControllerSet = wire.NewSet(
	AuthCtrlSet,
	UserCtrlSet,
	MenuCtrlSet,
	RoleCtrlSet,
	DigitPredictSet,
	AIModelCtrlSet,
)

var AuthCtrlSet = wire.NewSet(
	NewAuthController,
	wire.Struct(new(service.AuthService), "*"),
)

var UserCtrlSet = wire.NewSet(
	NewUserController,
	wire.Struct(new(service.UserService), "*"),
)

var MenuCtrlSet = wire.NewSet(
	NewMenuController,
	wire.Struct(new(service.MenuService), "*"),
)

var RoleCtrlSet = wire.NewSet(
	NewRoleController,
	wire.Struct(new(service.RoleService), "*"),
)

var DigitPredictSet = wire.NewSet(
	NewDigitPredictController,
	wire.Struct(new(service.DigitPredictService), "*"),
)

var AIModelCtrlSet = wire.NewSet(
	NewAIModelController,
	wire.Struct(new(service.AIModelService), "*"),
)
