package ai

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
)

// EmbeddingRequest 请求体
type EmbeddingRequest struct {
	EncodingFormat string   `json:"encoding_format"` // "float"
	Input          []string `json:"input"`           // 输入文本数组
	Model          string   `json:"model"`           // 模型名称
	Dimensions     int      `json:"dimensions"`      // 模型维度
}

// EmbeddingResponse 返回体结构
type EmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
		Object    string    `json:"object"`
	} `json:"data"`
	Model  string `json:"model"`
	Object string `json:"object"`
	Usage  struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// Embed 执行向量化请求
func (c *AIModelClient) Embed(ctx context.Context, apiKey string, input []string) ([][]float32, error) {
	// 直接指定用千问模型进行向量化
	reqBody := EmbeddingRequest{
		EncodingFormat: "float",
		Input:          input,
		Model:          "text-embedding-v4",
		Dimensions:     1024,
	}
	bodyBytes, _ := sonic.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost, "https://dashscope.aliyuncs.com/compatible-mode/v1/embeddings",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embedding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("embedding API error: %s", string(respData))
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var embeddingResp EmbeddingResponse
	if err := sonic.Unmarshal(respData, &embeddingResp); err != nil {
		return nil, fmt.Errorf("failed to decode response with sonic: %w", err)
	}

	vectors := make([][]float32, len(embeddingResp.Data))
	for _, item := range embeddingResp.Data {
		vectors[item.Index] = item.Embedding
	}

	return vectors, nil
}
