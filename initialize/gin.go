package initialize

import (
	"Art-Design-Backend/middleware"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func InitGin() *gin.Engine {
	// 初始化总路由
	r := gin.Default()
	// 通过Gzip压缩响应内容，减少传输数据量，提高传输速度。
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	// 添加操作日志记录
	r.Use(middleware.OperationLogger())
	// 添加校验器
	r.Use(middleware.ValidationErrorMiddleware())
	return r
}
