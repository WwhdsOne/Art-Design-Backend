package middleware

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/authutils"
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"time"
)

// OperationLoggerMiddleware 日志中间件
func (m *Middlewares) OperationLoggerMiddleware() gin.HandlerFunc {
	loggerChan := make(chan *entity.OperationLog, 1000)
	go func() {
		for log := range loggerChan {
			if err := m.Db.Create(log).Error; err != nil {
				zap.L().Error("保存操作日志失败", zap.Error(err))
			}
		}
	}()
	return func(c *gin.Context) {
		startTime := time.Now()
		var bodyBytes []byte
		// 检查 Content-Type 是否为 application/json
		contentType := c.GetHeader("Content-Type")
		if contentType == "application/json" {
			// 读取请求参数
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重置 body，供后续中间件使用
		}

		// 处理请求
		c.Next()

		// 获取用户信息 (示例代码，实际可通过 Token 或上下文获取)
		userID := authutils.GetUserID(c)

		// 收集响应信息
		statusCode := c.Writer.Status()

		// 创建操作日志
		log := &entity.OperationLog{
			OperatorID: userID,
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			IP:         c.ClientIP(),
			Params:     string(bodyBytes),
			Status:     int16(statusCode),
			CreatedAt:  startTime,
		}
		loggerChan <- log
	}
}
