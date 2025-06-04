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
		contentType := c.GetHeader("Content-Type")
		if contentType == "application/json" {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重置 body
		}

		// 执行请求
		c.Next()

		// 获取耗时（ms）
		latency := time.Since(startTime).Milliseconds()

		// 获取用户ID（根据你项目实际逻辑修改）
		userID := authutils.GetUserID(c)

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 获取 User-Agent
		userAgent := c.GetHeader("User-Agent")

		// 获取 Query 参数
		query := c.Request.URL.RawQuery

		var errMsg string
		errors := c.Errors
		if len(errors) > 0 {
			errMsg = errors[0].Error()
		}

		// 构造日志对象
		log := &entity.OperationLog{
			OperatorID: userID,
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			Query:      query,
			Params:     string(bodyBytes),
			Status:     int16(statusCode),
			Latency:    latency,
			IP:         c.ClientIP(),
			UserAgent:  userAgent,
			ErrorMsg:   errMsg,
			CreatedAt:  startTime,
		}

		loggerChan <- log
	}
}
