package response

import "github.com/dromara/carbon/v2"

// Role 定义角色模型
type Role struct {
	ID          int64           `json:"id,string"`
	Name        string          `json:"name"`
	Code        string          `json:"code"`
	Description string          `json:"description"`
	Status      int8            `json:"status"`
	CreatedAt   carbon.DateTime `json:"createdAt"`
}
