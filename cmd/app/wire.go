//go:build wireinject
// +build wireinject

package main

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/internal/bootstrap"
	"Art-Design-Backend/internal/controller"
	"Art-Design-Backend/internal/repository"

	"github.com/google/wire"
)

// 构造函数是因为初始化时有其他操作
// wire.Struct则只需要构造一个结构体
func wireApp() *bootstrap.HTTPServer {
	wire.Build(
		wire.Struct(new(bootstrap.HTTPServer), "*"),
		config.LoadConfig,
		config.ProvideDefaultUserConfig,
		config.ProviderMiddlewareConfig,
		bootstrap.InitSet,
		// 这里解释一下没有serviceProvider的原因:
		// 	service总是只被对应的controller使用，但是repo可能被多个service使用
		//  所以controllerProvider中直接就创建了service，没有单独的serviceProvider
		controller.ControllerSet,
		repository.RepositorySet,
	)
	return nil
}
