package middleware

import (
	"Art-Design-Backend/internal/model/entity"
	"Art-Design-Backend/pkg/authutils"
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mssola/useragent"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

type uaParsed struct {
	Browser  string
	Version  string
	OS       string
	Platform string
}

func (m *Middlewares) OperationLoggerMiddleware() gin.HandlerFunc {

	loggerChan := make(chan *entity.OperationLog, m.Config.OperationLog.LogChannelBufferSize)

	// UA 缓存：1小时过期，10分钟清理
	uaCache := cache.New(1*time.Hour, 10*time.Minute)

	// ================= 异步日志写入 =================
	go func() {
		for logItem := range loggerChan {

			uaStr := logItem.UserAgent
			var parsed *uaParsed

			if val, found := uaCache.Get(uaStr); found {
				parsed = val.(*uaParsed)
			} else {
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

			logItem.Browser = parsed.Browser
			logItem.BrowserVersion = parsed.Version
			logItem.OS = parsed.OS
			logItem.Platform = parsed.Platform

			if err := m.Db.Create(logItem).Error; err != nil {
				zap.L().Error("保存操作日志失败", zap.Error(err))
			}
		}
	}()

	// ================= Middleware =================
	return func(c *gin.Context) {

		startTime := time.Now()

		var bodyBytes []byte

		// ========= 仅必要时读取 body =========
		if shouldReadBody(c) {
			// 限制最大读取 1MB，防止大包 OOM
			limitedReader := io.LimitReader(c.Request.Body, 1<<20)
			bodyBytes, _ = io.ReadAll(limitedReader)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 执行后续 handler（包括 WS Upgrade）
		c.Next()

		// ========= 跳过 WS / SSE =========
		if isWebSocketRequest(c) || isSSERequest(c) {
			return
		}

		latency := time.Since(startTime).Milliseconds()
		userID := authutils.GetUserID(c)
		statusCode := c.Writer.Status()

		var errMsg string
		if len(c.Errors) > 0 {
			errMsg = c.Errors[0].Error()
		}

		log := &entity.OperationLog{
			OperatorID: userID,
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			Query:      c.Request.URL.RawQuery,
			Params:     string(bodyBytes),
			Status:     int16(statusCode),
			Latency:    latency,
			IP:         c.ClientIP(),
			UserAgent:  c.GetHeader("User-Agent"),
			ErrorMsg:   errMsg,
			CreatedAt:  startTime,
		}

		// 防止阻塞主流程
		select {
		case loggerChan <- log:
		default:
			zap.L().Warn("operation log channel is full, log dropped")
		}
	}
}

// ================= 判断函数 =================

func isWebSocketRequest(c *gin.Context) bool {
	return websocket.IsWebSocketUpgrade(c.Request)
}

func isSSERequest(c *gin.Context) bool {
	return strings.Contains(
		strings.ToLower(c.GetHeader("Accept")),
		"text/event-stream",
	)
}

func shouldReadBody(c *gin.Context) bool {

	if isWebSocketRequest(c) || isSSERequest(c) {
		return false
	}

	contentType := c.GetHeader("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		return false
	}

	switch c.Request.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return true
	default:
		return false
	}
}
