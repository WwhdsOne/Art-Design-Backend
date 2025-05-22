package response

import (
	"github.com/dromara/carbon/v2"
)

type DigitPredict struct {
	ID        int64           `json:"id,string"`
	Image     string          `json:"image"`
	Label     *int8           `json:"label"` // 预测结果
	CreatedAt carbon.DateTime `json:"created_at"`
}
