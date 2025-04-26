package global

import (
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// DB 运营平台数据库连接
	DB *gorm.DB
	// Logger 全局日志zapLog
	Logger *zap.Logger
	// OSSClient oss连接
	OSSClient *oss.Client
)
