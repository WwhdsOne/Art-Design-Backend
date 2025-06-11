package ai

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
)

// ChatRequest 普通非流式请求，返回完整响应体 []byte 或错误
func (c *AIModelClient) ChatRequest(ctx context.Context, url, token string, reqData ChatRequest) ([]byte, error) {
	url = fmt.Sprintf("%s/chat/completions", url)
	reqData.Stream = true
	body, err := sonic.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// ChatStream 流式请求，将远程 SSE 响应数据实时推送给 ginCtx，供前端实时消费
func (c *AIModelClient) ChatStream(ginCtx *gin.Context, url, token string, reqData ChatRequest) error {
	url = fmt.Sprintf("%s/chat/completions", url)
	reqData.Stream = true

	body, err := sonic.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ginCtx.Request.Context(), "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("bad response status: %d", resp.StatusCode)
	}

	// 设置 SSE 响应头
	ginCtx.Writer.Header().Set("Content-Type", "text/event-stream")
	ginCtx.Writer.Header().Set("Cache-Control", "no-cache")
	ginCtx.Writer.Header().Set("Connection", "keep-alive")
	ginCtx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ginCtx.Status(http.StatusOK)

	reader := bufio.NewReader(resp.Body)
	var responseContent string

	// 使用 Stream 自动 flush SSE 数据
	ginCtx.Stream(func(w io.Writer) bool {
		line, err := reader.ReadBytes('\n')
		zap.L().Info("Received line", zap.String("line", string(line)))
		if err != nil {
			_ = resp.Body.Close() // 确保清理
			return false          // 停止 stream
		}
		if len(bytes.TrimSpace(line)) == 0 {
			return true // 跳过空行
		}

		// 移除SSE前缀 "data: "
		line = bytes.TrimPrefix(line, []byte("data: "))

		// 跳过 "[DONE]" 标记
		if string(line) == "[DONE]" {
			_ = resp.Body.Close()
			return false
		}

		// 解析JSON响应
		var streamResponse ChatCompletionStreamResponse
		if err := sonic.Unmarshal(line, &streamResponse); err != nil {
			zap.L().Error("Failed to parse response", zap.Error(err), zap.String("raw", string(line)))
			return false
		}

		// 检查是否是结束标记
		if streamResponse.Choices[0].isEnd() {
			_ = resp.Body.Close()
			return false // 结束流式响应
		}

		// 提取content并添加到responseContent
		if content := streamResponse.Choices[0].Delta.Content; content != "" {
			// 构造 JSON 数据
			jsonData := map[string]string{"content": content}

			// 使用 sonic 序列化
			jsonBytes, err := sonic.Marshal(jsonData)
			if err != nil {
				zap.L().Error("JSON marshal error", zap.Error(err))
				return false
			}

			// 发送 JSON 格式的 SSE 数据
			if _, err := fmt.Fprintf(w, "data: %s\n\n", jsonBytes); err != nil {
				zap.L().Error("Write error", zap.Error(err))
				return false
			}

			// 刷新缓冲区
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}

		return true
	})

	// 返回完整的响应内容
	zap.L().Info("Complete response", zap.String("content", responseContent))
	return nil
}
