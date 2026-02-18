package ws

type Action struct {
	ActionID int64   `json:"action_id,string"`
	Action   string  `json:"action"`
	Selector *string `json:"selector,omitempty"`
	Value    *string `json:"value,omitempty"`
	URL      *string `json:"url,omitempty"`
	Distance *int    `json:"distance,omitempty"`
	Timeout  *int    `json:"timeout,omitempty"`
}

type PageElement struct {
	Tag      string  `json:"tag"`
	Text     string  `json:"text"`
	Selector string  `json:"selector"`
	Value    *string `json:"value,omitempty"`
	Visible  bool    `json:"visible"`  // 保留，前端不再传，默认 false
	Disabled bool    `json:"disabled"` // 保留，前端不再传，默认 false
}

type ScrollInfo struct {
	ScrollHeight float64 `json:"scrollHeight"` // 文档总高度
	ClientHeight float64 `json:"clientHeight"` // 可视区域高度
	ScrollTop    float64 `json:"scrollTop"`    // 当前滚动位置
	HasMoreBelow bool    `json:"hasMoreBelow"` // 下方是否还有内容
	HasMoreAbove bool    `json:"hasMoreAbove"` // 上方是否还有内容
}

type PageState struct {
	URL        string        `json:"url"`
	Title      string        `json:"title"`                // 页面标题
	Elements   []PageElement `json:"elements"`             // 页面元素列表
	ScrollInfo *ScrollInfo   `json:"scrollInfo,omitempty"` // 滚动信息，可选
}

type ClientMessage struct {
	Type          string     `json:"type"`
	MessageID     int64      `json:"message_id,string,omitempty"`
	PageState     *PageState `json:"pageState,omitempty"`
	ActionID      int64      `json:"action_id,string,omitempty"`
	Success       bool       `json:"success,omitempty"`
	Error         string     `json:"error,omitempty"`
	ExecutionTime int        `json:"execution_time,omitempty"`
	Task          string     `json:"task,omitempty"` // 新增
}

type ServerMessage struct {
	Type    string  `json:"type"`
	Action  *Action `json:"action,omitempty"`
	Message string  `json:"message,omitempty"`
}
