package ai

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"go.uber.org/zap"
	"io"
	"net/http"
)

// ChatRequest 普通非流式请求，返回完整响应体 []byte 或错误
func (c *AIModelClient) ChatRequest(ctx context.Context, url, token string, reqData ChatRequest) ([]byte, error) {
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

// ChatStreamWithWriter 流式请求，将远程 SSE 响应数据实时推送给 ginCtx，供前端实时消费
func (c *AIModelClient) ChatStreamWithWriter(ctx context.Context, w http.ResponseWriter, url, token string, reqData ChatRequest) (err error) {
	reqData.Stream = true

	body, err := sonic.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response status: %d", resp.StatusCode)
	}

	// 设置 SSE 响应头
	header := w.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
	header.Set("Access-Control-Allow-Origin", "*")

	// http.Flusher 必须支持刷新
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("response does not support flushing")
	}

	reader := bufio.NewReader(resp.Body)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var line []byte
			line, err = reader.ReadBytes('\n')
			if err != nil {
				return err
			}
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			line = bytes.TrimPrefix(line, []byte("data: "))
			if string(line) == "[DONE]" {
				return nil
			}

			var streamResponse ChatCompletionStreamResponse
			if err = sonic.Unmarshal(line, &streamResponse); err != nil {
				zap.L().Error("Failed to parse response", zap.Error(err), zap.String("raw", string(line)))
				return err
			}

			if streamResponse.Choices[0].isEnd() {
				return
			}

			if content := streamResponse.Choices[0].Delta.Content; content != "" {
				jsonData := map[string]string{"v": content}
				var jsonBytes []byte
				jsonBytes, err = sonic.Marshal(jsonData)
				if err != nil {
					zap.L().Error("JSON marshal error", zap.Error(err))
					return err
				}
				if _, err = fmt.Fprintf(w, "data: %s\n\n", jsonBytes); err != nil {
					zap.L().Error("Write error", zap.Error(err))
					return err
				}
				flusher.Flush()
			}
		}
	}
}
