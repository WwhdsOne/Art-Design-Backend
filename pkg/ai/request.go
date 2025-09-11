package ai

type MultiModeChatRequest struct {
	Messages         []MultiModeChatMessage `json:"messages"`
	Model            string                 `json:"model"`
	FrequencyPenalty *float64               `json:"frequency_penalty,omitempty"`
	MaxTokens        *int                   `json:"max_tokens,omitempty"`
	PresencePenalty  *float64               `json:"presence_penalty,omitempty"`
	ResponseFormat   *ResponseFormat        `json:"response_format,omitempty"`
	Stop             interface{}            `json:"stop,omitempty"` // string / []string / nil
	Stream           bool                   `json:"stream,omitempty"`
	StreamOptions    *StreamOptions         `json:"stream_options,omitempty"`
	Temperature      *float64               `json:"temperature,omitempty"`
	TopP             *float64               `json:"top_p,omitempty"`
	Tools            interface{}            `json:"tools,omitempty"`       // []Tool / nil
	ToolChoice       string                 `json:"tool_choice,omitempty"` // "none", "auto", "required"
	Logprobs         bool                   `json:"logprobs,omitempty"`
	TopLogprobs      interface{}            `json:"top_logprobs,omitempty"` // int / nil
}

type MultiModeChatMessage struct {
	Role    string                 `json:"role"`
	Content []MultiModeChatContent `json:"content"`
}

type MultiModeChatContent struct {
	Type     string `json:"type"`                // "text" or "image_url"
	Text     string `json:"text,omitempty"`      // 文本内容
	ImageUrl string `json:"image_url,omitempty"` // 图片 URL
}

type ChatMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type StreamOptions struct {
	// 根据具体 API 文档进行结构定义，这里暂为占位
}

type ChatRequest struct {
	Messages         []ChatMessage   `json:"messages"`
	Model            string          `json:"model"`
	FrequencyPenalty *float64        `json:"frequency_penalty,omitempty"`
	MaxTokens        *int            `json:"max_tokens,omitempty"`
	PresencePenalty  *float64        `json:"presence_penalty,omitempty"`
	ResponseFormat   *ResponseFormat `json:"response_format,omitempty"`
	Stop             interface{}     `json:"stop,omitempty"` // string / []string / nil
	Stream           bool            `json:"stream,omitempty"`
	StreamOptions    *StreamOptions  `json:"stream_options,omitempty"`
	Temperature      *float64        `json:"temperature,omitempty"`
	TopP             *float64        `json:"top_p,omitempty"`
	Tools            interface{}     `json:"tools,omitempty"`       // []Tool / nil
	ToolChoice       string          `json:"tool_choice,omitempty"` // "none", "auto", "required"
	Logprobs         bool            `json:"logprobs,omitempty"`
	TopLogprobs      interface{}     `json:"top_logprobs,omitempty"` // int / nil
}

func DefaultStreamChatRequest(model string, messages []ChatMessage) ChatRequest {
	return ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}
}

func DefaultChatRequest(model string, messages []ChatMessage) ChatRequest {
	return ChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}
}

func DefaultMultiModeChatRequest(model string, messages []MultiModeChatMessage) MultiModeChatRequest {
	return MultiModeChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}
}
