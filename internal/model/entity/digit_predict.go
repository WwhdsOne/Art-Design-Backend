package entity

import (
	"Art-Design-Backend/internal/model/base"
	"Art-Design-Backend/pkg/constant/tablename"
)

type DigitPredict struct {
	base.BaseModel
	Image string `gorm:"type:varchar(255);comment:图片地址"`
	Label *int8  `gorm:"type:tinyint;comment:识别结果"`
}

func (d *DigitPredict) TableName() string {
	return tablename.DigitPredict
}
