package response

import "time"

type ConversationResponse struct {
	ID        int64     `json:"id,string"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MessageResponse struct {
	ID             int64     `json:"id,string"`
	ConversationID int64     `json:"conversation_id,string"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

type ActionResponse struct {
	ID            int64     `json:"id,string"`
	MessageID     int64     `json:"message_id,string"`
	ActionType    string    `json:"action_type"`
	Sequence      int       `json:"sequence"`
	Status        string    `json:"status"`
	URL           *string   `json:"url,omitempty"`
	Selector      *string   `json:"selector,omitempty"`
	Value         *string   `json:"value,omitempty"`
	Distance      *int      `json:"distance,omitempty"`
	Timeout       *int      `json:"timeout,omitempty"`
	ErrorMessage  *string   `json:"error_message,omitempty"`
	ExecutionTime *int      `json:"execution_time,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type MessageWithActionsResponse struct {
	Message MessageResponse  `json:"message"`
	Actions []ActionResponse `json:"actions"`
}
