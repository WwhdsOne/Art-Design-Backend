package ai

// ChatCompletionStreamResponse 表示流式聊天完成API的响应
type ChatCompletionStreamResponse struct {
	ID                string         `json:"id"`
	Object            string         `json:"object"`
	Created           int64          `json:"created"`
	Model             string         `json:"model"`
	SystemFingerprint string         `json:"system_fingerprint,omitempty"`
	Choices           []StreamChoice `json:"choices"`
}

// StreamChoice 表示流式响应中的一个选择
type StreamChoice struct {
	Index        int          `json:"index"`
	Delta        DeltaContent `json:"delta"`
	Logprobs     interface{}  `json:"logprobs,omitempty"`
	FinishReason *string      `json:"finish_reason,omitempty"`
}

// DeltaContent 表示流式响应中的增量内容
type DeltaContent struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// isEnd 判断是否是流式响应的最后一个数据块
func (c *StreamChoice) isEnd() bool {
	// 检查FinishReason
	if c.FinishReason != nil && *c.FinishReason != "" {
		return true
	}

	// 检查内容是否为[DONE]
	return c.Delta.Content == "[DONE]"
}
