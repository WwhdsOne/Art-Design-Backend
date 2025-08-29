package ai

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
	"go.uber.org/zap"
)

// ChatRequest 普通非流式请求，返回完整响应体 []byte 或错误
func (c *AIModelClient) ChatRequest(ctx context.Context, url, token string, reqData ChatRequest) ([]byte, error) {
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

// ChatStreamWithWriter 流式请求，将远程 SSE 响应数据实时推送给 ginCtx，
// 同时在函数返回时拼接完整的 AI 响应
func (c *AIModelClient) ChatStreamWithWriter(
	ctx context.Context,
	w http.ResponseWriter,
	url, token string,
	reqData ChatRequest,
) (fullResp string, err error) {
	reqData.Stream = true

	body, err := sonic.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad response status: %d", resp.StatusCode)
	}

	// 设置 SSE 响应头
	header := w.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Cache-Control", "no-cache")
	header.Set("Connection", "keep-alive")
	header.Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return "", fmt.Errorf("response does not support flushing")
	}

	reader := bufio.NewReader(resp.Body)

	var sb strings.Builder // 用于拼接完整响应

	for {
		select {
		case <-ctx.Done():
			return sb.String(), ctx.Err()
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil {
				return sb.String(), err
			}
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			line = bytes.TrimPrefix(line, []byte("data: "))
			if string(line) == "[DONE]" {
				return sb.String(), nil
			}

			var streamResponse ChatCompletionStreamResponse
			if err = sonic.Unmarshal(line, &streamResponse); err != nil {
				zap.L().Error("Failed to parse response", zap.Error(err), zap.String("raw", string(line)))
				return sb.String(), err
			}

			if streamResponse.Choices[0].isEnd() {
				return sb.String(), nil
			}

			if content := streamResponse.Choices[0].Delta.Content; content != "" {
				sb.WriteString(content) // 拼接到完整响应

				// 实时推送到前端
				jsonData := map[string]string{"v": content}
				jsonBytes, _ := sonic.Marshal(jsonData)
				if _, err = fmt.Fprintf(w, "data: %s\n\n", jsonBytes); err != nil {
					return sb.String(), err
				}
				flusher.Flush()
			}
		}
	}
}
