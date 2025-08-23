package middleware

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/authutils"
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

// uaParsed 缓存结构
type uaParsed struct {
	Browser  string
	Version  string
	OS       string
	Platform string
}

// OperationLoggerMiddleware 日志中间件（带 UA 缓存）
func (m *Middlewares) OperationLoggerMiddleware() gin.HandlerFunc {
	loggerChan := make(chan *entity.OperationLog, m.Config.OperationLog.LogChannelBufferSize)

	// 使用 go-cache：ttl = 永久（不主动过期），清理间隔也可以设长一点
	uaCache := cache.New(1*time.Hour, 10*time.Minute)

	// 异步日志写入逻辑
	go func() {
		for logItem := range loggerChan {
			uaStr := logItem.UserAgent
			var parsed *uaParsed

			// 优先读取缓存
			if val, found := uaCache.Get(uaStr); found {
				parsed = val.(*uaParsed)
			} else {
				// 解析 UA 并缓存
				ua := useragent.New(uaStr)
				browser, version := ua.Browser()
				parsed = &uaParsed{
					Browser:  browser,
					Version:  version,
					OS:       ua.OS(),
					Platform: ua.Platform(),
				}
				uaCache.Set(uaStr, parsed, cache.DefaultExpiration)
			}

			// 写入解析信息
			logItem.Browser = parsed.Browser
			logItem.BrowserVersion = parsed.Version
			logItem.OS = parsed.OS
			logItem.Platform = parsed.Platform

			// 写入数据库
			if err := m.Db.Create(logItem).Error; err != nil {
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
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重新填充请求体
		}

		// 继续处理请求
		c.Next()

		// 获取耗时（ms）
		latency := time.Since(startTime).Milliseconds()

		// 获取用户ID（根据你项目实际逻辑修改）
		userID := authutils.GetUserID(c)

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 获取 User-Agent
		userAgent := c.GetHeader("User-Agent")

		// 获取 URL 参数
		query := c.Request.URL.RawQuery

		// 错误信息
		var errMsg string
		if len(c.Errors) > 0 {
			errMsg = c.Errors[0].Error()
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
