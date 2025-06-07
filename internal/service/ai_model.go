package service

import (
	"Art-Design-Backend/internal/repository"
	"Art-Design-Backend/pkg/redisx"
)

type AIModelService struct {
	AIModelRepo *repository.AIModelRepository      // AI模型
	GormTX      *repository.GormTransactionManager // 事务
	Redis       *redisx.RedisWrapper               // redis
}
