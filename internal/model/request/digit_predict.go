package request

import "Art-Design-Backend/internal/model/base"

type DigitPredict struct {
	ID    base.LongStringID `json:"id" binding:"required"`
	Image string            `json:"image" binding:"required"`
}
