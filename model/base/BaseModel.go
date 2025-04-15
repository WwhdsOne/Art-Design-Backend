package base

import (
	"Art-Design-Backend/pkg/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

// BaseModel 是一个通用的模型，包含 ID、CreatedAt、UpdatedAt 和 DeletedAt 字段
type BaseModel struct {
	ID        int64     `gorm:"type:bigint;column:id;primarykey"`                // 雪花ID
	CreatedAt time.Time `gorm:"type:timestamp;column:created_at;autoCreateTime"` // 创建时间字段，记录记录创建时的时间戳
	UpdatedAt time.Time `gorm:"type:timestamp;column:updated_at;autoUpdateTime"` // 更新时间字段，记录最后一次更新时的时间戳
	UpdateBy  int64     `gorm:"type:bigint;column:updated_by"`                   // 修改人字段，记录最后一次更新操作者的标识
	CreateBy  int64     `gorm:"type:bigint;column:created_by"`                   // 创建人字段，记录创建操作者的标识
}

// BeforeCreate 在创建记录之前自动生成雪花 ID
func (b *BaseModel) BeforeCreate(db *gorm.DB) (err error) {
	b.ID = utils.GenerateSnowflakeId()
	b.fillAddReq(db.Statement.Context)
	return
}

func (b *BaseModel) BeforeUpdate(db *gorm.DB) (err error) {
	b.fillUpdateReq(db.Statement.Context)
	return
}

func (b *BaseModel) fillAddReq(c context.Context) {
	userID := utils.GetUserID(c)
	// 如果存在 claims，正常提取用户 ID
	b.CreateBy = userID
	b.UpdateBy = userID
}

func (b *BaseModel) fillUpdateReq(c context.Context) {
	userID := utils.GetUserID(c)
	b.UpdateBy = userID
}
