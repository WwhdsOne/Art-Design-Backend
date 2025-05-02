package config

import (
	"Art-Design-Backend/internal/controller"
	"Art-Design-Backend/pkg/middleware"
	"Art-Design-Backend/pkg/utils"
	"context"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	App          string `yaml:"app"`
	Port         string `yaml:"port"`
	ReadTimeOut  string `yaml:"read-time-out"`
	WriteTimeOut string `yaml:"write-time-out"`
}

type HttpServer struct {
	Engine         *gin.Engine                // gin引擎
	Logger         *zap.Logger                // 日志
	AuthController *controller.AuthController // 鉴权Ctrl
	Config         *Config                    // 服务器配置
}

func NewGin(m *middleware.Middlewares, logger *zap.Logger) *gin.Engine {
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
	// 添加操作日志记录
	return engine
}

func (h *HttpServer) GinServer() {
	cfg := h.Config
	httpServer := http.Server{
		Addr:         cfg.Server.Port,
		Handler:      h.Engine,
		ReadTimeout:  utils.ParseDuration(cfg.Server.ReadTimeOut),
		WriteTimeout: utils.ParseDuration(cfg.Server.WriteTimeOut),
	}
	// 协程启动http服务
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// 设置信号通道
	stopChan := make(chan os.Signal, 1)
	// 监听SIGINT和SIGTERM信号
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞直到收到信号
	<-stopChan
	h.Logger.Info("正在关闭服务器...")

	// 创建一个上下文用于优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 尝试优雅关闭服务器
	if err := httpServer.Shutdown(ctx); err != nil {
		h.Logger.Fatal("服务器关闭失败: ", zap.Error(err))
	}
	h.Logger.Info("服务器已正常退出")
}
