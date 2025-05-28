package bootstrap

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/pkg/middleware"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"regexp"
	"time"
)

// RegisterValidator 注册全局请求校验器
func RegisterValidator() {
	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("strongpassword", func(fl validator.FieldLevel) bool {
			pass := fl.Field().String()
			hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pass)
			hasLower := regexp.MustCompile(`[a-z]`).MatchString(pass)
			hasNumber := regexp.MustCompile(`[0-9]`).MatchString(pass)
			return hasUpper && hasLower && hasNumber
		})
		if err != nil {
			zap.L().Fatal("自定义校验器注册失败")
			return
		}
	}
}

func InitGin(m *middleware.Middlewares, logger *zap.Logger, c *config.Config) *gin.Engine {
	// 注册自定义校验器
	RegisterValidator()
	// 创建gin引擎
	engine := gin.New()
	// 通过Gzip压缩响应内容，减少传输数据量，提高传输速度。
	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	// 设置日志
	engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(logger, true))
	// 设置全局错误校验器错误中间件
	engine.Use(m.ErrorHandlerMiddleware())
	// 设置操作日志数据库记录中间件
	engine.Use(m.OperationLoggerMiddleware())
	// 添加限流中间件
	engine.Use(m.RedisRateLimitMiddleware(c.Server.RateLimit.Window, c.Server.RateLimit.MaxReq))
	// 添加操作日志记录
	return engine
}
