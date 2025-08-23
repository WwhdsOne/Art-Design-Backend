package common

import (
	"Art-Design-Backend/pkg/authutils"
	"time"

	"gorm.io/gorm"
)

// BaseModel 是一个通用的模型，包含 ID、CreatedAt、UpdatedAt 和 DeletedAt 字段
type BaseModel struct {
	ID        int64     `gorm:"type:bigserial;column:id;primaryKey"`
	CreatedAt time.Time `gorm:"type:timestamp;column:created_at;autoCreateTime"` // 创建时间字段，记录记录创建时的时间戳
	UpdatedAt time.Time `gorm:"type:timestamp;column:updated_at;autoUpdateTime"` // 更新时间字段，记录最后一次更新时的时间戳
	UpdateBy  int64     `gorm:"type:bigint;column:updated_by"`                   // 修改人字段，记录最后一次更新操作者的标识
	CreateBy  int64     `gorm:"type:bigint;column:created_by"`                   // 创建人字段，记录创建操作者的标识
}

func (b *BaseModel) BeforeCreate(db *gorm.DB) (err error) {
	operatorID := authutils.GetUserID(db.Statement.Context)
	b.CreateBy = operatorID
	b.UpdateBy = operatorID
	return
}

func (b *BaseModel) BeforeUpdate(db *gorm.DB) (err error) {
	operatorID := authutils.GetUserID(db.Statement.Context)
	b.UpdateBy = operatorID
	return
}
