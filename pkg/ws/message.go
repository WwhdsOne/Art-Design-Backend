package ws

type Action struct {
	ActionID     int64   `json:"action_id,string"`
	Action       string  `json:"action"`
	Index        *int    `json:"index,omitempty"`
	OptionIndex  *int    `json:"option_index,omitempty"`  // 用于选择题的选项索引
	Selector     *string `json:"selector,omitempty"`
	Value        *string `json:"value,omitempty"`
	URL          *string `json:"url,omitempty"`
	Distance     *int    `json:"distance,omitempty"`
	Timeout      *int    `json:"timeout,omitempty"`
}

type SelectOption struct {
	Value string `json:"value"`
	Text  string `json:"text"`
}

type PageElement struct {
	Index       int             `json:"index"`
	Tag         string          `json:"tag"`
	Text        string          `json:"text"`
	Selector    string          `json:"selector"`
	Type        *string         `json:"type,omitempty"`           // 元素类型（input, button, question等）
	Question    *string         `json:"question,omitempty"`      // 题目文本（仅当type=question时有值）
	Value       *string         `json:"value,omitempty"`
	Label       *string         `json:"label,omitempty"`
	Role        *string         `json:"role,omitempty"`
	AriaLabel   *string         `json:"ariaLabel,omitempty"`
	AriaExpanded *string        `json:"ariaExpanded,omitempty"`
	AriaChecked  *string        `json:"ariaChecked,omitempty"`
	Required    *bool           `json:"required,omitempty"`
	Disabled    *bool           `json:"disabled,omitempty"`
	Options     []PageElement   `json:"options,omitempty"`      // 选项列表（仅question类型使用）
	IsNew       *bool           `json:"isNew,omitempty"`
}

type ScrollInfo struct {
	ScrollHeight float64 `json:"scrollHeight"`
	ClientHeight float64 `json:"clientHeight"`
	ScrollTop    float64 `json:"scrollTop"`
	HasMoreBelow bool    `json:"hasMoreBelow"`
	HasMoreAbove bool    `json:"hasMoreAbove"`
}

type PageState struct {
	URL         string        `json:"url"`
	Title       string        `json:"title"`
	Elements    []PageElement `json:"elements"`
	ElementText string       `json:"elementText,omitempty"`
	ScrollInfo  *ScrollInfo   `json:"scrollInfo,omitempty"`
	Screenshot  string       `json:"screenshot,omitempty"` // base64 编码的带标签截图
}

// AgentOutput 结构化思维输出（LLM 返回的 JSON）
type AgentOutput struct {
	Thinking     string   `json:"thinking,omitempty"`
	Evaluation   string   `json:"evaluation_previous_goal,omitempty"`
	Memory       string   `json:"memory,omitempty"`
	NextGoal     string   `json:"next_goal,omitempty"`
	Action       string   `json:"action,omitempty"`
	Actions      []Action `json:"actions,omitempty"`
	NeedVision   bool     `json:"need_vision,omitempty"` // 文本模型请求视觉模型辅助
}

type ClientMessage struct {
	Type          string     `json:"type"`
	MessageID     int64      `json:"message_id,string,omitempty"`
	PageState     *PageState `json:"pageState,omitempty"`
	ActionID      int64      `json:"action_id,string,omitempty"`
	Success       bool       `json:"success,omitempty"`
	Error         string     `json:"error,omitempty"`
	ExecutionTime int        `json:"execution_time,omitempty"`
	Task          string     `json:"task,omitempty"`
	// 多动作结果
	Results       []ActionResult `json:"results,omitempty"`
}

type ActionResult struct {
	ActionID      string `json:"action_id"`
	Success       bool   `json:"success"`
	ExecutionTime int    `json:"execution_time"`
	Error         string `json:"error,omitempty"`
}

type ServerMessage struct {
	Type             string    `json:"type"`
	Action           *Action   `json:"action,omitempty"`
	Actions          []Action  `json:"actions,omitempty"`
	StopOnPageChange *bool     `json:"stop_on_page_change,omitempty"`
	Message          string    `json:"message,omitempty"`
}
