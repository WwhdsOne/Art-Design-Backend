//go:build wireinject
// +build wireinject

package main

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/internal/controller"
	"Art-Design-Backend/pkg/middleware"
	"github.com/google/wire"
)

// 构造函数是因为初始化时有其他操作
// wire.Struct则只需要构造一个结构体
func wireApp() *config.HttpServer {
	wire.Build(
		wire.Struct(new(config.HttpServer), "*"),
		config.NewConfig,
		config.NewLogger,
		config.NewRedis,
		config.NewGorm,
		wire.Struct(new(middleware.Middlewares), "*"),
		config.NewGin,
		config.NewJWT,
		controller.ControllersProvider,
	)
	return nil
}
