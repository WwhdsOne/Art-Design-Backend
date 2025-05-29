package request

import "Art-Design-Backend/internal/model/base"

type ChangeStatus struct {
	ID     base.LongStringID `json:"id" binding:"required"`
	Status int8              `json:"status" binding:"required"`
}
