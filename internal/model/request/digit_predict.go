package request

import "Art-Design-Backend/internal/model/common"

type DigitPredict struct {
	ID    common.LongStringID `json:"id" label:"ID" binding:"required"`
	Image string              `json:"image" label:"图片URL" binding:"required"`
}
