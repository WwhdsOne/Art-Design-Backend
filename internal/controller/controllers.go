package controller

import (
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/internal/service"
	"github.com/google/wire"
)

var ControllersProvider = wire.NewSet(AuthCtrlProvider)

var AuthCtrlProvider = wire.NewSet(
	NewAuthController,
	repository.NewUserRepository,
	wire.Struct(new(service.AuthService), "*"),
)
