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
		cCp := c.Copy()
		// 保存日志到数据库
		go func() { // 异步保存，避免阻塞请求
			if err := m.Db.WithContext(cCp).Create(log).Error; err != nil {
				// 打印日志或记录到其他地方
				zap.L().Error("Failed to save operation log\n")
				// 打印具体错误
				zap.L().Error(err.Error())
			}
		}()
	}
}
