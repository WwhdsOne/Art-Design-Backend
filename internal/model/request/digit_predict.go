package request

import "Art-Design-Backend/internal/model/base"

type DigitPredict struct {
	ID    base.LongStringID `json:"id" label:"ID" binding:"required"`
	Image string            `json:"image" label:"图片URL" binding:"required"`
}
