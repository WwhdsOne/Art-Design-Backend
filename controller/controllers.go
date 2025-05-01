package controller

import (
	"Art-Design-Backend/service"
	"github.com/google/wire"
)

var ControllersProvider = wire.NewSet(AuthCtrlProvider)

var AuthCtrlProvider = wire.NewSet(
	NewAuthController,
	wire.Struct(new(service.AuthService), "*"),
)
