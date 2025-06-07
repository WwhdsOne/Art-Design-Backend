package middleware

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/authutils"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"go.uber.org/zap"
	"io"
	"time"
)

// uaParsed 缓存结构
type uaParsed struct {
	Browser  string
	Version  string
	OS       string
	Platform string
}

const maxUACacheSize = 1000

// OperationLoggerMiddleware 日志中间件（带 UA 缓存）
func (m *Middlewares) OperationLoggerMiddleware() gin.HandlerFunc {
	loggerChan := make(chan *entity.OperationLog, 1000)

	// UA 缓存（map + FIFO keys）
	uaCache := make(map[string]*uaParsed)
	uaKeys := make([]string, 0, maxUACacheSize)

	// 后台异步处理日志写入
	go func() {
		for log := range loggerChan {
			uaStr := log.UserAgent
			var parsed *uaParsed

			if val, ok := uaCache[uaStr]; ok {
				parsed = val
			} else {
				ua := useragent.New(uaStr)
				browser, version := ua.Browser()
				parsed = &uaParsed{
					Browser:  browser,
					Version:  version,
					OS:       ua.OS(),
					Platform: ua.Platform(),
				}

				// 缓存维护
				if len(uaKeys) >= maxUACacheSize {
					oldest := uaKeys[0]
					uaKeys = uaKeys[1:]
					delete(uaCache, oldest)
				}
				uaKeys = append(uaKeys, uaStr)
				uaCache[uaStr] = parsed
			}

			// 写入解析结果
			log.Browser = parsed.Browser
			log.BrowserVersion = parsed.Version
			log.OS = parsed.OS
			log.Platform = parsed.Platform

			// 写入数据库
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
