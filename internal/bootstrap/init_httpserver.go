package bootstrap

import (
	"Art-Design-Backend/config"
	"Art-Design-Backend/internal/controller"
	"Art-Design-Backend/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer struct {
	Engine                 *gin.Engine                        // gin引擎
	Logger                 *zap.Logger                        // 日志
	AuthController         *controller.AuthController         // 鉴权Ctrl
	UserController         *controller.UserController         // 用户Ctrl
	MenuController         *controller.MenuController         // 菜单Ctrl
	RoleController         *controller.RoleController         // 角色Ctrl
	DigitPredictController *controller.DigitPredictController // 数字预测Ctrl
	AIController           *controller.AIController           // AI模型Ctrl
	Config                 *config.Config                     // 服务器配置
}

func (h *HttpServer) InitGinServer() {
	cfg := h.Config
	httpServer := http.Server{
		Addr:         cfg.Server.Port,
		Handler:      h.Engine,
		ReadTimeout:  utils.ParseDuration(cfg.Server.ReadTimeOut),
		WriteTimeout: utils.ParseDuration(cfg.Server.WriteTimeOut),
		IdleTimeout:  utils.ParseDuration(cfg.Server.IdleTimeout),
	}
	// 协程启动http服务
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			h.Logger.Fatal("服务器启动失败:", zap.Error(err))
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
	c, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 尝试优雅关闭服务器
	if err := httpServer.Shutdown(c); err != nil {
		h.Logger.Fatal("服务器关闭失败: ", zap.Error(err))
	}
	h.Logger.Info("服务器已正常退出")
}
