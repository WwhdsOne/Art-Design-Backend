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
	Logprobs     any          `json:"logprobs,omitempty"`
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

// ChatCompletionResponse 非流式返回的标准结构体
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"` // 时间戳
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   *ChatCompletionUsage   `json:"usage,omitempty"`
}

// ChatCompletionChoice 每条生成结果
type ChatCompletionChoice struct {
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

// ChatCompletionMessage 消息结构体
type ChatCompletionMessage struct {
	Role    string `json:"role"`    // "user" / "assistant" / "system"
	Content string `json:"content"` // 生成的文本
}

// ChatCompletionUsage token 使用情况
type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// FirstText 可选：辅助方法，快速获取第一条生成文本
func (r *ChatCompletionResponse) FirstText() string {
	if len(r.Choices) > 0 {
		return r.Choices[0].Message.Content
	}
	return ""
}
