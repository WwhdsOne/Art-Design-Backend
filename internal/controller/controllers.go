package controller

import (
	"Art-Design-Backend/internal/service"
	"github.com/google/wire"
)

var ControllersProvider = wire.NewSet(AuthCtrlProvider, UserCtrlProvider)

var AuthCtrlProvider = wire.NewSet(
	NewAuthController,
	wire.Struct(new(service.AuthService), "*"),
)

var UserCtrlProvider = wire.NewSet(
	NewUserController,
	wire.Struct(new(service.UserService), "*"),
)
