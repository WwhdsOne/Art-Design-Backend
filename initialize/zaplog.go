package initialize

import (
	"go.uber.org/zap"
)

// InitLogger 初始化 Zap Logger
func InitLogger() (logger *zap.Logger) {
	production, err := zap.NewProduction()
	if err != nil {
		return nil
	}
	return production
}
