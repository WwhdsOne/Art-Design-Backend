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

type Position struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type PageElement struct {
	Tag      string    `json:"tag"`
	Text     string    `json:"text"`
	Selector string    `json:"selector"`
	Value    *string   `json:"value,omitempty"`
	Type     *string   `json:"type,omitempty"`
	Label    *string   `json:"label,omitempty"`
	Position *Position `json:"position,omitempty"`
}

type ScrollInfo struct {
	ScrollHeight float64 `json:"scrollHeight"`
	ClientHeight float64 `json:"clientHeight"`
	ScrollTop    float64 `json:"scrollTop"`
	HasMoreBelow bool    `json:"hasMoreBelow"`
	HasMoreAbove bool    `json:"hasMoreAbove"`
}

type PageState struct {
	URL        string        `json:"url"`
	Title      string        `json:"title"`
	Elements   []PageElement `json:"elements"`
	ScrollInfo *ScrollInfo   `json:"scrollInfo,omitempty"`
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
}

type ServerMessage struct {
	Type    string  `json:"type"`
	Action  *Action `json:"action,omitempty"`
	Message string  `json:"message,omitempty"`
}
