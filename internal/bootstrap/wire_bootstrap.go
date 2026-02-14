package bootstrap

import "github.com/google/wire"

var InitSet = wire.NewSet(
	InitLogger,
	InitRedis,
	InitGorm,
	InitMiddleware,
	InitDigitPredict,
	InitAIModelClient,
	InitGin,
	InitOSSClient,
	InitJWT,
	InitSlicer,
	InitWebSocketHub,
)
