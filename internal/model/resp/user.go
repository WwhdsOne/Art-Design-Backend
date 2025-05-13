package resp

import (
	"github.com/dromara/carbon/v2"
)

type User struct {
	ID           int64           `json:"id,string"`
	RealName     string          `json:"realname"`
	Nickname     string          `json:"nickname"`
	Gender       int8            `json:"gender"`
	Email        string          `json:"email"`
	Phone        string          `json:"phone"`
	Address      string          `json:"address"`
	Avatar       string          `json:"avatar"`
	Introduction string          `json:"introduction"`
	Occupation   string          `json:"occupation"`
	Tags         []string        `json:"tags"`
	Roles        []Role          `json:"roles"`
	Status       int8            `json:"status"`
	CreatedAt    carbon.DateTime `json:"created_at"`
}
