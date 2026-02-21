package ai

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/bytedance/sonic"
)

// RerankRequest 保持原有结构体定义不变
type RerankRequest struct {
	Model           string   `json:"model"`
	Query           string   `json:"query"`
	Documents       []string `json:"documents"`
	TopN            *int     `json:"top_n,omitempty"`
	ReturnDocuments *bool    `json:"return_documents,omitempty"`
	MaxChunksPerDoc *int     `json:"max_chunks_per_doc,omitempty"`
	OverlapTokens   *int     `json:"overlap_tokens,omitempty"`
}

// RerankResponse结构体根据实际响应格式进行调整
type RerankResponse struct {
	ID      string       `json:"id"`
	Results []ResultItem `json:"results"`
	Meta    MetaInfo     `json:"meta"`
}

// ResultItem表示每个排序结果项的结构体
type ResultItem struct {
	Index          int     `json:"index"`
	RelevanceScore float64 `json:"relevance_score"`
}

// MetaInfo包含元数据相关信息的结构体
type MetaInfo struct {
	BilledUnits BilledUnits `json:"billed_units"`
	Tokens      Tokens      `json:"tokens"`
}

// BilledUnits记录计费相关的各类token数量等信息
type BilledUnits struct {
	InputTokens     int `json:"input_tokens"`
	OutputTokens    int `json:"output_tokens"`
	SearchUnits     int `json:"search_units"`
	Classifications int `json:"classifications"`
}

// Tokens记录输入输出token相关信息
type Tokens struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// Rerank 使用sonic进行序列化的重排序调用方法
func (c *AIModelClient) Rerank(token string, req RerankRequest, topK int) ([]string, error) {
	reqBody, err := sonic.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		http.MethodPost,
		"https://api.siliconflow.cn/v1/rerank",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send rerank request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// 读取响应体以获取更多错误详情
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("bad response status: %d, failed to read response body: %w", resp.StatusCode, readErr)
		}

		return nil, fmt.Errorf("bad response status: %d, response body: %s", resp.StatusCode, string(respBody))
	}

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result RerankResponse
	if err := sonic.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	if topK > 0 && len(result.Results) > topK {
		result.Results = result.Results[:topK]
	}

	// 根据得分高低从req的Documents中获取前topK个文档内容
	type ResultWithDoc struct {
		Score float64
		Doc   string
	}
	resultsWithDoc := make([]ResultWithDoc, len(result.Results))
	for i, r := range result.Results {
		resultsWithDoc[i].Score = r.RelevanceScore
		resultsWithDoc[i].Doc = req.Documents[r.Index]
	}

	// 对结果按照得分从高到低排序
	sort.Slice(resultsWithDoc, func(i, j int) bool {
		return resultsWithDoc[i].Score > resultsWithDoc[j].Score
	})

	topKDocuments := make([]string, topK)
	for i := range topK {
		topKDocuments[i] = resultsWithDoc[i].Doc
	}

	return topKDocuments, nil
}
