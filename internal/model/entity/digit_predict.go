package entity

import (
	"Art-Design-Backend/internal/model/common"
	"Art-Design-Backend/pkg/constant/tablename"
)

type DigitPredict struct {
	common.BaseModel
	Image string `gorm:"type:varchar(255);comment:图片地址"`
	Label *int8  `gorm:"type:smallint;comment:识别结果"`
}

func (d *DigitPredict) TableName() string {
	return tablename.DigitPredict
}
