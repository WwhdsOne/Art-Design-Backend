package initialize

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化 Zap Logger
// InitLogger 初始化日志记录器，根据条件调整日志级别
func InitLogger(shouldSuppressDebug bool) (logger *zap.Logger) {
	config := zap.NewDevelopmentConfig()
	if shouldSuppressDebug {
		// 如果满足条件，将日志级别设置为 InfoLevel
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	logger, _ = config.Build()
	return
}
