package request

import "Art-Design-Backend/internal/model/common"

type ChangeStatus struct {
	ID     common.LongStringID `json:"id" binding:"required"`
	Status int8                `json:"status" binding:"required"`
}
