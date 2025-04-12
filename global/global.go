package global

import (
	"Art-Design-Backend/pkg/jwt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// DB 运营平台数据库连接
	DB *gorm.DB
	// Redis Redis连接
	Redis *redis.Client
	// JWT 全局JWT
	JWT *jwt.JWT
	// Logger 全局日志zapLog
	Logger *zap.Logger
	// OSSClient oss连接
	OSSClient *oss.Client
)
