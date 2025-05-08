package resp

// Role 定义角色模型
type Role struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status int8   `json:"status"`
}
