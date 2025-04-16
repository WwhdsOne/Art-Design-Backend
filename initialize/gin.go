package initialize

import (
	"Art-Design-Backend/global"
	"Art-Design-Backend/pkg/middleware"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"time"
)

func InitGin() *gin.Engine {
	// 初始化总路由
	r := gin.New()
	// 通过Gzip压缩响应内容，减少传输数据量，提高传输速度。
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	// 设置日志
	r.Use(ginzap.Ginzap(global.Logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(global.Logger, true))
	// 添加操作日志记录
	r.Use(middleware.OperationLogger())
	// 添加全局错误校验器错误中间件
	r.Use(middleware.ErrorHandlingMiddleware())
	return r
}
