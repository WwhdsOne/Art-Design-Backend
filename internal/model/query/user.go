package query

import (
	"Art-Design-Backend/internal/model/common"
)

type User struct {
	RealName string `json:"realname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Gender   int8   `json:"gender"`
	Status   int8   `json:"status"`
	common.PaginationReq
}
