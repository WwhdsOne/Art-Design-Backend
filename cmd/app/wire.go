//go:build wireinject
// +build wireinject

package main

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/internal/controller"
	"Art-Design-Backend/internal/repository"
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
		config.NewOSSClient,
		// 这里解释一下没有serviceProvider的原因:
		// 	service总是只被对应的controller使用，但是repo可能被多个service使用
		//  所以controllerProvider中直接就创建了service，没有单独的serviceProvider
		controller.ControllersProvider,
		repository.RepositoriesProvider,
	)
	return nil
}
