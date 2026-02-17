package response

import "time"

type ConversationResponse struct {
	ID          int64     `json:"id,string"`
	Title       string    `json:"title"`
	State       string    `json:"state"`
	BrowserType string    `json:"browser_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MessageResponse struct {
	ID             int64     `json:"id,string"`
	ConversationID int64     `json:"conversation_id,string"`
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

type AdminSummaryResponse struct {
	TodayTasks    int64  `json:"todayTasks"`
	TodayGrowth   string `json:"todayGrowth"`
	SuccessRate   int    `json:"successRate"`
	SuccessGrowth string `json:"successGrowth"`
}

type VolumeDataResponse struct {
	Volume    int64   `json:"volume"`
	Growth    string  `json:"growth"`
	ChartData []int64 `json:"chartData"`
}

type RateDataResponse struct {
	Rate      int     `json:"rate"`
	Growth    string  `json:"growth"`
	ChartData []int64 `json:"chartData"`
}

type TotalTaskVolumeResponse struct {
	Total        int64              `json:"total"`
	TotalAll     int64              `json:"totalAll"`
	Distribution []DistributionItem `json:"distribution"`
}

type DistributionItem struct {
	Value int64  `json:"value"`
	Name  string `json:"name"`
}

type TaskClassificationResponse struct {
	Distribution []ClassificationItem `json:"distribution"`
}

type ClassificationItem struct {
	Value int    `json:"value"`
	Name  string `json:"name"`
}

type ActiveSessionsResponse struct {
	Count     int64   `json:"count"`
	Growth    string  `json:"growth"`
	ChartData []int64 `json:"chartData"`
}

type AnnualTaskStatsResponse struct {
	Year        int     `json:"year"`
	MonthlyData []int64 `json:"monthlyData"`
	QuarterData []int64 `json:"quarterData"`
}

type HotTaskItemResponse struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	AvgTime     int    `json:"avgTime"`
	ExecCount   int64  `json:"execCount"`
	SuccessRate int    `json:"successRate"`
	Color       string `json:"color"`
}

type UserSummaryResponse struct {
	SessionCount  int64  `json:"sessionCount"`
	SessionGrowth string `json:"sessionGrowth"`
	SuccessRate   int    `json:"successRate"`
	SuccessGrowth string `json:"successGrowth"`
}

type UserTaskOverviewResponse struct {
	SessionCount int64         `json:"sessionCount"`
	TaskCount    int64         `json:"taskCount"`
	SuccessRate  int           `json:"successRate"`
	WeekGrowth   string        `json:"weekGrowth"`
	ChartData    UserTaskChart `json:"chartData"`
}

type UserTaskChart struct {
	XAxis []string `json:"xAxis"`
	Data  []int64  `json:"data"`
}

type UserTaskTrendResponse struct {
	Year        int     `json:"year"`
	Growth      string  `json:"growth"`
	MonthlyData []int64 `json:"monthlyData"`
}
