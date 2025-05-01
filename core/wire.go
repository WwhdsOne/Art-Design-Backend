//go:build wireinject
// +build wireinject

package core

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/controller"
	"Art-Design-Backend/pkg/middleware"
	"github.com/google/wire"
)

func wireApp() *config.HttpServer {
	wire.Build(
		wire.Struct(new(config.HttpServer), "*"),
		config.NewConfig,
		config.NewLogger,
		config.NewRedis,
		config.NewGorm,
		middleware.NewMiddlewares,
		config.NewGin,
		config.NewJWT,
		controller.ControllersProvider,
	)
	return nil
}
