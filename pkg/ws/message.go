package ws

type Action struct {
	Action   string  `json:"action"`
	Selector *string `json:"selector,omitempty"`
	Value    *string `json:"value,omitempty"`
	URL      *string `json:"url,omitempty"`
	Distance *int    `json:"distance,omitempty"`
	Timeout  *int    `json:"timeout,omitempty"`
}

type PageElement struct {
	Tag      string `json:"tag"`
	Text     string `json:"text"`
	Selector string `json:"selector"`
	Visible  bool   `json:"visible"`
	Disabled bool   `json:"disabled"`
}

type PageState struct {
	URL      string        `json:"url"`
	Elements []PageElement `json:"elements"`
}

type ClientMessage struct {
	Type      string     `json:"type"`
	Task      string     `json:"task,omitempty"`
	PageState *PageState `json:"pageState,omitempty"`
	Success   bool       `json:"success,omitempty"`
	Error     string     `json:"error,omitempty"`
}

type ServerMessage struct {
	Type    string  `json:"type"`
	Action  *Action `json:"action,omitempty"`
	Message string  `json:"message,omitempty"`
}
