package request

import "Art-Design-Backend/internal/model/base"

type Role struct {
	ID          base.LongStringID `json:"id"`
	Name        string            `json:"name" binding:"required,max=10"`
	Code        string            `json:"code" binding:"required,min=2,max=10"`
	Description string            `json:"description"`
	Status      int8              `json:"status"`
}
