package base

import (
	"Art-Design-Backend/pkg/utils"
	"database/sql/driver"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// ID 是一个通用的 ID 类型，包含一个 int64 值
type ID struct {
	Val int64 `gorm:"column:id;type:bigint;primarykey" json:"id,string"`
}

func (i *ID) BeforeCreate(db *gorm.DB) (err error) {
	// 单条记录生成 ID
	id, err := utils.GenerateSnowflakeId()
	if err == nil {
		i.Val = id
		return
	}
	return
}

func (i *ID) Value() (driver.Value, error) {
	return i.Val, nil
}

func (i *ID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("cannot scan nil value into ID")
	}
	switch v := value.(type) {
	case int64:
		i.Val = v
		return nil
	default:
		return fmt.Errorf("cannot scan value of type %T into ID", value)
	}
}
