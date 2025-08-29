package response

import "time"

type Conversation struct {
	ID        int64     `json:"id,string"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
