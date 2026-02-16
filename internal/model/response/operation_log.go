package response

import "time"

type OperationLog struct {
	ID             int64     `json:"id"`
	OperatorID     int64     `json:"operator_iD"`
	Method         string    `json:"method"`
	Path           string    `json:"path"`
	Query          string    `json:"query"`
	Params         string    `json:"params"`
	Status         int16     `json:"status"`
	Latency        int64     `json:"latency"`
	IP             string    `json:"ip"`
	UserAgent      string    `json:"user_agent"`
	Browser        string    `json:"browser"`
	BrowserVersion string    `json:"browser_version"`
	OS             string    `json:"os"`
	Platform       string    `json:"platform"`
	ErrorMsg       string    `json:"error_msg"`
	CreatedAt      time.Time `json:"created_at"`
}
